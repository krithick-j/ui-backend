CREATE SCHEMA IF NOT EXISTS ui_network;

USE ui_network;

CREATE TABLE IF NOT EXISTS `bank_names` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_bank_names_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE IF NOT EXISTS `bv_transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `disrib_id` varchar(30) DEFAULT NULL,
  `order_id` varchar(30) DEFAULT NULL,
  `date` datetime(3) DEFAULT NULL,
  `bv_value` bigint(20) DEFAULT NULL,
  `is_active` int(1) NOT NULL DEFAULT 0,
  `activate_date` datetime(3) DEFAULT NULL,
  `place` varchar(30) DEFAULT NULL,
  `side` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_bv_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `cart_items` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `quantity` bigint(20) unsigned DEFAULT 1,
  `distrib_id` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cart_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=595 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `cheque_frequencies` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `frequency` bigint(20) DEFAULT NULL,
  `place` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cheque_frequencies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `contact_us` (
  `name` varchar(30) NOT NULL,
  `distrib_id` varchar(30) NOT NULL,
  `country` varchar(30) NOT NULL,
  `contact_number` varchar(30) NOT NULL,
  `email_address` varchar(30) NOT NULL,
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `admin_name` varchar(30) DEFAULT NULL,
  `aadhaar_front` varchar(30) DEFAULT NULL,
  `aadhaar_back` varchar(30) DEFAULT NULL,
  `pan_card` varchar(30) DEFAULT NULL,
  `passport_size` varchar(30) DEFAULT NULL,
  `enquiry_type_id` bigint(20) unsigned NOT NULL,
  `your_query` varchar(30) DEFAULT NULL,
  `pan_number` varchar(30) DEFAULT NULL,
  `account_holder_name` varchar(30) DEFAULT NULL,
  `account_number` varchar(30) DEFAULT NULL,
  `bank_name_id` bigint(20) unsigned DEFAULT NULL,
  `bank_ifsc_code` varchar(30) DEFAULT NULL,
  `bank_branch` varchar(30) DEFAULT NULL,
  `refund_type` bigint(20) DEFAULT NULL,
  `orders_header_id` bigint(20) unsigned DEFAULT NULL,
  `status` tinyint(1) DEFAULT 1,
  PRIMARY KEY (`id`),
  KEY `idx_contact_us_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `enquiry_types` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `admin_name` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_enquiry_types_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `i_coupon_transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `v_id` varchar(30) DEFAULT NULL,
  `pin` varchar(30) DEFAULT NULL,
  `total_value` double DEFAULT NULL,
  `amount_detected` double DEFAULT NULL,
  `balance` double DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  `order_id` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_i_coupon_transactions_deleted_at` (`deleted_at`),
  KEY `fk_i_coupon_transactions_product` (`product_id`),
  CONSTRAINT `fk_i_coupon_transactions_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=158 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `i_coupons` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `date_on` datetime(3) DEFAULT NULL,
  `tx_detail` varchar(30) DEFAULT NULL,
  `admin_name` varchar(30) DEFAULT NULL,
  `v_id` varchar(30) DEFAULT NULL,
  `rank` decimal(10,2) DEFAULT NULL,
  `value` double DEFAULT NULL,
  `expires_on` datetime(3) DEFAULT NULL,
  `pin` varchar(30) DEFAULT NULL,
  `active` tinyint(1) DEFAULT NULL,
  `balance` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_i_coupons_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=9640 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `orders_headers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `order_id` varchar(30) DEFAULT NULL,
  `sub_total` double DEFAULT NULL,
  `total_sand_h` double DEFAULT NULL,
  `total_amount` double DEFAULT NULL,
  `total_quantity` double DEFAULT NULL,
  `total_bv` bigint(20) DEFAULT NULL,
  `contact_name` varchar(30) DEFAULT NULL,
  `contact_email` varchar(30) DEFAULT NULL,
  `address` varchar(30) DEFAULT NULL,
  `city` varchar(30) DEFAULT NULL,
  `district` varchar(30) DEFAULT NULL,
  `state` varchar(30) DEFAULT NULL,
  `zip_code` bigint(20) unsigned DEFAULT NULL,
  `country` varchar(30) DEFAULT NULL,
  `home_phone_no` varchar(30) DEFAULT NULL,
  `mobile_phone_no` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_headers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=167 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `orders_liners` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `orders_header_id` bigint(20) unsigned DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `quantity` bigint(20) unsigned DEFAULT NULL,
  `unit_price` bigint(20) unsigned DEFAULT NULL,
  `bv` bigint(20) DEFAULT NULL,
  `sub_total` double DEFAULT NULL,
  `sand_h` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_liners_deleted_at` (`deleted_at`),
  KEY `fk_orders_headers_orders_liner` (`orders_header_id`),
  CONSTRAINT `fk_orders_headers_orders_liner` FOREIGN KEY (`orders_header_id`) REFERENCES `orders_headers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=237 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `product_categories` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `product_images` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `image` varchar(30) DEFAULT NULL,
  `product_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_images_deleted_at` (`deleted_at`),
  KEY `fk_products_product_images` (`product_id`),
  CONSTRAINT `fk_products_product_images` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `quantity` bigint(20) DEFAULT NULL,
  `shipment_time` varchar(30) DEFAULT NULL,
  `price` double DEFAULT NULL,
  `sand_h` double DEFAULT NULL,
  `rsp` bigint(20) DEFAULT NULL,
  `product_category_id` bigint(20) DEFAULT NULL,
  `bv` bigint(20) DEFAULT NULL,
  `ep` double DEFAULT NULL,
  `admin_name` varchar(30) DEFAULT NULL,
  `product_image_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `fk_products_product_category` (`product_category_id`),
  CONSTRAINT `fk_products_product_category` FOREIGN KEY (`product_category_id`) REFERENCES `product_categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `rsp_transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `order_id` varchar(30) DEFAULT NULL,
  `rsp` float DEFAULT NULL,
  `date` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rsp_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tracking_centers` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `created_at` datetime(3) DEFAULT NULL,
  `distrib_id` varchar(30) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `place` varchar(30) DEFAULT NULL,
  `p_place` varchar(30) DEFAULT NULL,
  `p_distrib_id` varchar(30) DEFAULT NULL,
  `left_distrib_id` varchar(30) DEFAULT NULL,
  `left_place` varchar(30) DEFAULT NULL,
  `right_distrib_id` varchar(30) DEFAULT NULL,
  `right_place` varchar(30) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bv` bigint(20) DEFAULT NULL,
  `left_point` bigint(20) DEFAULT 0,
  `right_point` bigint(20) DEFAULT 0,
  `is_active` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tracking_centers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=243 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `user_rank` (
  `distrib_id` varchar(100) NOT NULL,
  `year` int(10) unsigned NOT NULL,
  `month` int(10) unsigned NOT NULL,
  `rank` decimal(8,2) DEFAULT NULL,
  PRIMARY KEY (`distrib_id`,`year`,`month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `distrib_id` varchar(30) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `pass` varchar(30) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `address1` varchar(30) DEFAULT NULL,
  `address2` varchar(30) DEFAULT NULL,
  `town_or_city` varchar(30) DEFAULT NULL,
  `district` varchar(30) DEFAULT NULL,
  `state_or_province` varchar(30) DEFAULT NULL,
  `email_address` varchar(30) DEFAULT NULL,
  `pin_or_zip_code` bigint(20) unsigned DEFAULT NULL,
  `country` varchar(30) DEFAULT NULL,
  `home_phone_no` varchar(30) DEFAULT NULL,
  `mobile_phone_no` varchar(30) DEFAULT NULL,
  `ref_distrib_id` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='USER TABLE TO BE MODIFIED LATER';
