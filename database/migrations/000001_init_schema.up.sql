-- 1. Table: users
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('barista', 'customer') NOT NULL DEFAULT 'customer',
    points INT NOT NULL DEFAULT 0,
    fcm_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE vouchers (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    type ENUM('menu_promo', 'cart_discount') NOT NULL DEFAULT 'cart_discount',
    discount_amount INT NOT NULL,
    min_purchase INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 3. Table: products (Terkait dengan Vouchers)
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    discount INT NOT NULL DEFAULT 0,
    voucher_id VARCHAR(36), -- FK ke tabel vouchers (Nullable)
    volume VARCHAR(50) NOT NULL,
    image_urls JSON,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE SET NULL
);

-- 4. Table: ingredients (Terkait dengan Products)
CREATE TABLE ingredients (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    grammage VARCHAR(50) NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- 5. Table: raw_materials (Berdiri Sendiri)
CREATE TABLE raw_materials (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL, -- Mendukung angka desimal (misal: 1.5)
    unit VARCHAR(50) NOT NULL,
    price INT NOT NULL,
    source VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 6. Table: orders (Terkait dengan Users)
CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36), -- Nullable untuk tamu/non-user
    total_amount INT NOT NULL,
    discount INT NOT NULL DEFAULT 0,
    final_amount INT NOT NULL,
    is_redeem BOOLEAN NOT NULL DEFAULT FALSE,
    status ENUM('pending', 'completed', 'cancelled') NOT NULL DEFAULT 'completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- Tidak ada updated_at karena order bersifat final (immutable record)
);

-- 7. Table: order_items (Terkait dengan Orders dan Products)
CREATE TABLE order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL,
    price_at_time INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- 8. Table: point_transactions (Terkait dengan Users dan Orders)
CREATE TABLE point_transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36), -- Nullable karena saat redeem token, order_id mungkin menyusul/terpisah
    type ENUM('earned', 'redeemed') NOT NULL,
    points INT NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
);

-- 9. Table: system_settings
CREATE TABLE system_settings (
    setting_key VARCHAR(100) PRIMARY KEY,
    setting_value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert Data Awal (Opsional)
INSERT INTO system_settings (setting_key, setting_value) 
VALUES ('barista_status', 'offline');