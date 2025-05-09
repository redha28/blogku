-- Insert admin user
INSERT INTO admins (username, password, email, created_at, updated_at) VALUES
('admin', '$argon2id$v=19$m=65536,t=3,p=2$Zusw4bvnBRM9i4LkwHTBtw$65J1ViiiJCCkLP1XjwfvwmgbvZWf6srIPx9ZVZoQq2k', 'admin@blog.com', '2025-05-08 16:19:31', '2025-05-08 16:19:31');
