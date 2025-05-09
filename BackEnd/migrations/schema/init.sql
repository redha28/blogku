-- Create database if not exists
CREATE DATABASE IF NOT EXISTS blogku;
USE blogku;

-- Create tables
CREATE TABLE IF NOT EXISTS admins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blogs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    image_path VARCHAR(255),
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert admin user
INSERT INTO admins (username, email, password) VALUES 
('admin', 'admin@example.com', '$2a$10$zNybpz7ViULF7lhUs7VUH.MIrKzIcnuEhi5eCS2jfcTLdq0JlpL4m');

-- Insert sample blogs
INSERT INTO blogs (title, content, slug, image_path, published_at) VALUES
('Samsung Galaxy S24 Ultra Rilis, Bawa Kamera 200MP!','Samsung resmi merilis Galaxy S24 Ultra sebagai flagship terbarunya untuk tahun ini. Salah satu fitur unggulan dari perangkat ini adalah kamera utama beresolusi 200MP yang mampu mengambil gambar dengan detail luar biasa, bahkan dalam kondisi cahaya rendah.\n\nDitenagai oleh prosesor Snapdragon 8 Gen 3, Galaxy S24 Ultra menghadirkan performa tinggi untuk gaming, multitasking, dan penggunaan intensif lainnya. Layar 6.8 inci Dynamic AMOLED 2X dengan refresh rate 120Hz juga memastikan tampilan yang mulus dan jernih.\n\nBaterai berkapasitas 5000mAh diklaim mampu bertahan seharian penuh, dengan dukungan pengisian cepat 45W. Tak ketinggalan, perangkat ini hadir dengan fitur-fitur AI terbaru dari Samsung, serta dukungan S-Pen seperti pada generasi sebelumnya.\n\nGalaxy S24 Ultra menjalankan One UI 6 berbasis Android 14 dan akan menerima update sistem operasi hingga 4 tahun ke depan. Smartphone ini menjadi pilihan utama bagi pengguna yang menginginkan performa dan fitur kamera terbaik di kelasnya.','samsung-galaxy-s24-ultra-rilis-bawa-kamera-200mp','samsung-galaxy-s24-ultra-rilis-bawa-kamera-200mp.jpg',NOW()),
('iPhone 15 Pro Max Hadir dengan Titanium dan USB-C','Apple akhirnya menghadirkan perubahan besar pada iPhone 15 Pro Max dengan penggunaan bahan titanium yang membuatnya lebih ringan namun tetap kokoh. Ini merupakan langkah besar dari Apple yang sebelumnya menggunakan stainless steel.\n\nTak hanya itu, iPhone 15 Pro Max juga hadir dengan port USB-C menggantikan Lightning, memungkinkan kecepatan transfer data yang lebih tinggi dan kompatibilitas yang lebih luas dengan aksesori modern. Prosesor A17 Pro yang dibangun dengan arsitektur 3nm memberikan efisiensi daya yang lebih baik dan kinerja grafis yang meningkat, terutama untuk gaming berat dan aplikasi berbasis AI.\n\nSistem kamera juga mengalami peningkatan, dengan lensa periskop yang memungkinkan zoom optik hingga 5x tanpa kehilangan detail. Fitur ProRAW dan ProRes juga ditingkatkan, membuat iPhone 15 Pro Max menjadi perangkat ideal bagi fotografer dan videografer profesional.\n\nDengan iOS 17, Apple menyematkan fitur-fitur seperti StandBy Mode, Journal App, dan peningkatan privasi yang lebih ketat. iPhone 15 Pro Max benar-benar menggabungkan desain elegan dengan performa luar biasa.','iphone-15-pro-max-hadir-dengan-titanium-dan-usb-c','iphone-15-pro-max-hadir-dengan-titanium-dan-usb-c_image.webp',NOW());
