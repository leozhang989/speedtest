/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50641
 Source Host           : localhost:3306
 Source Schema         : speedtest

 Target Server Type    : MySQL
 Target Server Version : 50641
 File Encoding         : 65001

 Date: 30/08/2019 18:56:01
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations`  (
  `id_migration` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'surrogate key',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'migration name, unique',
  `created_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'date migrated or rolled back',
  `statements` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'SQL statements for this migration',
  `rollback_statements` longtext CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT 'SQL statment for rolling back migration',
  `status` enum('update','rollback') CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'update indicates it is a normal migration while rollback means this migration is rolled back',
  PRIMARY KEY (`id_migration`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, 'Settings_20190830_175200', '2019-08-30 17:52:03', 'CREATE TABLE settings(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT \'ID\',`setting_key` varchar(255) NOT NULL DEFAULT \'\',`setting_value` varchar(255) NOT NULL DEFAULT \'\',PRIMARY KEY (`id`))', NULL, 'update');
INSERT INTO `migrations` VALUES (2, 'Users_20190830_183829', '2019-08-30 18:38:32', 'CREATE TABLE users(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT \'ID\',`device_code` varchar(100) NOT NULL DEFAULT \'\',`vip_expiration_time` int(11) NOT NULL DEFAULT 0,`original_transaction_id` varchar(30) NOT NULL DEFAULT \'\',`updated` int(11) NOT NULL DEFAULT 0,`created` int(11) NOT NULL DEFAULT 0,PRIMARY KEY (`id`))', NULL, 'update');
INSERT INTO `migrations` VALUES (3, 'Orders_20190830_184415', '2019-08-30 18:44:19', 'CREATE TABLE orders(`id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT \'ID\',`device_code` varchar(100) NOT NULL DEFAULT \'\',`pay_status` tinyint(1) NOT NULL,`certificate` longtext NOT NULL,`latest_receipt` longtext NOT NULL,`updated` int(11) NOT NULL DEFAULT 0,`created` int(11) NOT NULL DEFAULT 0,PRIMARY KEY (`id`))', NULL, 'update');

-- ----------------------------
-- Table structure for orders
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `device_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `pay_status` tinyint(1) NOT NULL,
  `certificate` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `latest_receipt` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `updated` int(11) NOT NULL DEFAULT 0,
  `created` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for settings
-- ----------------------------
DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `setting_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `setting_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `device_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `vip_expiration_time` int(11) NOT NULL DEFAULT 0,
  `original_transaction_id` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `updated` int(11) NOT NULL DEFAULT 0,
  `created` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
