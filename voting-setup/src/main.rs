use clap::Parser;
use libsums::{client::SumsClient, member};
use mongodb::{
    bson::{doc, Document},
    options::ClientOptions,
    Client,
};

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    /// The MongoDB username
    mongo_user: String,

    /// The MongoDB password
    mongo_pass: String,

    /// The MongoDB URL. Note that this isn't a full URL (e.g., mongodb://...), but is just the
    /// address name.
    mongo_url: String,

    /// The Webdriver address for libsums
    webdriver_address: String,

    /// The browser being used for libsums (probably chrome)
    browser_name: String,

    /// The username to login to SUMS with. Must be a valid HackSoc commitee member.
    sums_username: String,

    /// The password for the SUMS user.
    sums_password: String,
}

#[tokio::main]
async fn main() {
    let args = Args::parse();

    let client_options = ClientOptions::parse(format!(
        "mongodb://{}:{}@{}/hacksoc?authSource=admin&retryWrites=true&w=majority",
        args.mongo_user, args.mongo_pass, args.mongo_url,
    ))
    .await
    .expect("Faileed to parse client options!");

    let mongodb_client =
        Client::with_options(client_options).expect("Failed to create MongoDB client!");

    let members_collection = mongodb_client
        .database("Hacksoc")
        .collection::<Document>("members");

    let sums_client = SumsClient::new(213, "http://localhost:4444", "chrome")
        .await
        .expect("Failed to create SUMS client!");

    println!("Authenticating to SUMS...");

    sums_client
        .authenticate(args.sums_username, args.sums_password)
        .await
        .expect("Failed to authenticate with SUMS!");

    println!("Getting member list...");

    let members = sums_client.members().await.expect("Failed to get members!");
    let student_ids = members
        .iter()
        .filter_map(|member| member.student_id.parse::<u32>().ok());

    let member_bson = student_ids.map(|student_id| doc! { "ID": student_id });

    members_collection
        .drop(None)
        .await
        .expect("Failed to drop members collection!");

    members_collection
        .insert_many(member_bson, None)
        .await
        .expect("Failed to insert members!");

    println!("Successfully inserted {} members", members.len());
}
