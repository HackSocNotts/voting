db.auth('admin-user', 'admin-password')

db = db.getSiblingDB('test-database')

db.createUser({
  user: 'brucebernard',
  pwd: 'admin',
  roles: [
    {
      role: 'root',
      db: 'hacksoc',
    },
  ],
});