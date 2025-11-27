-- MySQL dump 10.13  Distrib 8.0.19, for Win64 (x86_64)
--
-- Host: localhost    Database: TTDB
-- ------------------------------------------------------
-- Server version	8.4.7

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `branch`
--

DROP TABLE IF EXISTS `branch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branch` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` json DEFAULT NULL,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `branch`
--

LOCK TABLES `branch` WRITE;
/*!40000 ALTER TABLE `branch` DISABLE KEYS */;
INSERT INTO `branch` VALUES (1,'{\"en\": \"Bangkok Branch 1 (Updated)\", \"th\": \"สาขา กทม 1 (อัปเดตแล้ว)\"}','2025-11-25 08:05:47'),(2,'{\"en\": \"Bangkok Branch 2\", \"th\": \"สาขา กทม 2\"}','2025-11-19 06:47:03'),(3,'{\"en\": \"Chiang Mai Branch\", \"th\": \"สาขา เชียงใหม่\"}','2025-11-19 06:47:03'),(4,'{\"en\": \"Phuket Branch\", \"th\": \"สาขา ภูเก็ต\"}','2025-11-19 06:47:03'),(5,'{\"en\": \"Khon Kaen Branch\", \"th\": \"สาขา ขอนแก่น\"}','2025-11-19 06:47:03'),(6,'{\"en\": \"Nonthaburi Branch\", \"th\": \"สาขา นนทบุรี\"}','2025-11-19 06:47:03'),(7,'{\"en\": \"Nakhon Ratchasima Branch\", \"th\": \"สาขา นครราชสีมา\"}','2025-11-19 06:47:03'),(8,'{\"en\": \"Chonburi Branch\", \"th\": \"สาขา ชลบุรี\"}','2025-11-19 06:47:03'),(9,'{\"en\": \"Surat Thani Branch\", \"th\": \"สาขา สุราษฎร์ธานี\"}','2025-11-19 06:47:03'),(10,'{\"en\": \"Ubon Ratchathani Branch\", \"th\": \"สาขา อุบลราชธานี\"}','2025-11-19 06:47:03'),(13,'{\"en\": \"Bangkok Branch 1 (Updated)\", \"th\": \"สาขา กทม 1 (อัปเดตแล้ว)\"}','2025-11-21 09:26:35');
/*!40000 ALTER TABLE `branch` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `branch_location`
--

DROP TABLE IF EXISTS `branch_location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branch_location` (
  `id` int NOT NULL AUTO_INCREMENT,
  `branch_id` int NOT NULL,
  `province_id` int NOT NULL,
  `name` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_branch_location_branch` (`branch_id`),
  CONSTRAINT `fk_branch_location_branch` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `branch_location`
--

LOCK TABLES `branch_location` WRITE;
/*!40000 ALTER TABLE `branch_location` DISABLE KEYS */;
INSERT INTO `branch_location` VALUES (1,1,10,NULL),(2,2,20,NULL),(3,3,30,NULL),(4,4,40,NULL),(5,5,50,NULL),(6,6,60,NULL),(7,7,70,NULL),(8,8,80,NULL),(9,9,90,NULL),(10,10,100,NULL);
/*!40000 ALTER TABLE `branch_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `branches_interests`
--

DROP TABLE IF EXISTS `branches_interests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branches_interests` (
  `branch_id` int NOT NULL,
  `interest_id` int NOT NULL,
  PRIMARY KEY (`branch_id`,`interest_id`),
  KEY `fk_branch_interest_interest` (`interest_id`),
  CONSTRAINT `fk_branch_interest_branch` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_branch_interest_interest` FOREIGN KEY (`interest_id`) REFERENCES `interest` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `branches_interests`
--

LOCK TABLES `branches_interests` WRITE;
/*!40000 ALTER TABLE `branches_interests` DISABLE KEYS */;
INSERT INTO `branches_interests` VALUES (1,1),(9,1),(10,1),(1,2),(2,2),(10,2),(1,3),(2,3),(3,3),(2,4),(3,4),(4,4),(3,5),(4,5),(5,5),(4,6),(5,6),(6,6),(5,7),(6,7),(7,7),(6,8),(7,8),(8,8),(7,9),(8,9),(9,9),(8,10),(9,10),(10,10);
/*!40000 ALTER TABLE `branches_interests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `branches_products`
--

DROP TABLE IF EXISTS `branches_products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branches_products` (
  `branch_id` int NOT NULL,
  `product_id` int NOT NULL,
  PRIMARY KEY (`branch_id`,`product_id`),
  KEY `fk_branch_product_product` (`product_id`),
  CONSTRAINT `fk_branch_product_branch` FOREIGN KEY (`branch_id`) REFERENCES `branch` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_branch_product_product` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `branches_products`
--

LOCK TABLES `branches_products` WRITE;
/*!40000 ALTER TABLE `branches_products` DISABLE KEYS */;
INSERT INTO `branches_products` VALUES (2,1),(2,2),(3,2),(3,3),(4,3),(4,4),(5,4),(1,5),(5,5),(6,5),(13,5),(1,6),(6,6),(7,6),(13,6),(1,7),(7,7),(8,7),(13,7),(8,8),(9,8),(13,8),(9,9),(10,9),(10,10);
/*!40000 ALTER TABLE `branches_products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `interest`
--

DROP TABLE IF EXISTS `interest`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `interest` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` json DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `interest`
--

LOCK TABLES `interest` WRITE;
/*!40000 ALTER TABLE `interest` DISABLE KEYS */;
INSERT INTO `interest` VALUES (1,'{\"en\": \"Sports\", \"th\": \"กีฬา\"}'),(2,'{\"en\": \"Technology\", \"th\": \"เทคโนโลยี\"}'),(3,'{\"en\": \"Travel\", \"th\": \"ท่องเที่ยว\"}'),(4,'{\"en\": \"Food\", \"th\": \"อาหาร\"}'),(5,'{\"en\": \"Health\", \"th\": \"สุขภาพ\"}'),(6,'{\"en\": \"Music\", \"th\": \"ดนตรี\"}'),(7,'{\"en\": \"Books\", \"th\": \"หนังสือ\"}'),(8,'{\"en\": \"Games\", \"th\": \"เกม\"}'),(9,'{\"en\": \"Fashion\", \"th\": \"แฟชั่น\"}'),(10,'{\"en\": \"Pets\", \"th\": \"สัตว์เลี้ยง\"}');
/*!40000 ALTER TABLE `interest` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `outbox_events`
--

DROP TABLE IF EXISTS `outbox_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `outbox_events` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `aggregate_id` varchar(255) NOT NULL,
  `aggregate_type` varchar(255) NOT NULL,
  `event_type` varchar(50) NOT NULL,
  `payload` json DEFAULT NULL,
  `status` enum('pending','processed','failed') NOT NULL DEFAULT 'pending',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_created_at` (`status`,`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `outbox_events`
--

LOCK TABLES `outbox_events` WRITE;
/*!40000 ALTER TABLE `outbox_events` DISABLE KEYS */;
INSERT INTO `outbox_events` VALUES (1,'1','branch','updated','{\"id\": 1, \"name\": {\"en\": \"Bangkok Branch 1 (Updated)\", \"th\": \"สาขา กทม 1 (อัปเดตแล้ว)\"}, \"product_ids\": [5, 6, 7]}','processed','2025-11-25 08:05:47');
/*!40000 ALTER TABLE `outbox_events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` json DEFAULT NULL,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product`
--

LOCK TABLES `product` WRITE;
/*!40000 ALTER TABLE `product` DISABLE KEYS */;
INSERT INTO `product` VALUES (1,'{\"en\": \"Product A\", \"th\": \"สินค้า A\"}','2025-11-19 06:44:18'),(2,'{\"en\": \"Product B\", \"th\": \"สินค้า B\"}','2025-11-19 06:44:18'),(3,'{\"en\": \"Product C\", \"th\": \"สินค้า C\"}','2025-11-19 06:44:18'),(4,'{\"en\": \"Product D\", \"th\": \"สินค้า D\"}','2025-11-19 06:44:18'),(5,'{\"en\": \"Product E\", \"th\": \"สินค้า E\"}','2025-11-19 06:44:18'),(6,'{\"en\": \"Product F\", \"th\": \"สินค้า F\"}','2025-11-19 06:44:18'),(7,'{\"en\": \"Product G\", \"th\": \"สินค้า G\"}','2025-11-19 06:44:18'),(8,'{\"en\": \"Product H\", \"th\": \"สินค้า H\"}','2025-11-19 06:44:18'),(9,'{\"en\": \"Product I\", \"th\": \"สินค้า I\"}','2025-11-19 06:44:18'),(10,'{\"en\": \"Product J\", \"th\": \"สินค้า J\"}','2025-11-19 06:44:18');
/*!40000 ALTER TABLE `product` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_option`
--

DROP TABLE IF EXISTS `product_option`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_option` (
  `id` int NOT NULL AUTO_INCREMENT,
  `normal_price_thb` float DEFAULT NULL,
  `tagthai_price_thb` float DEFAULT NULL,
  `product_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_product_option_product` (`product_id`),
  CONSTRAINT `fk_product_option_product` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_option`
--

LOCK TABLES `product_option` WRITE;
/*!40000 ALTER TABLE `product_option` DISABLE KEYS */;
INSERT INTO `product_option` VALUES (1,100,90,1),(2,150,120,1),(3,200,180,2),(4,250,210,2),(5,300,270,3),(6,350,300,3),(7,400,360,4),(8,450,390,4),(9,500,450,5),(10,550,480,5),(11,600,540,6),(12,650,570,6),(13,700,630,7),(14,750,660,7),(15,800,720,8),(16,850,750,8),(17,900,810,9),(18,950,840,9),(19,1000,900,10),(20,1100,990,10);
/*!40000 ALTER TABLE `product_option` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'TTDB'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-11-25 15:21:24
