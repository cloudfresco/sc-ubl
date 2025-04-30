/*!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.6.18-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: scfdb1
-- ------------------------------------------------------
-- Server version	10.6.18-MariaDB-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `addresses`
--

DROP TABLE IF EXISTS `addresses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `addresses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `name1` varchar(50) DEFAULT '',
  `addr_list_agency_id` varchar(50) DEFAULT '',
  `addr_list_id` varchar(50) DEFAULT '',
  `addr_list_version_id` varchar(50) DEFAULT '',
  `address_type_code` varchar(10) DEFAULT '',
  `address_format_code` varchar(10) DEFAULT '',
  `postbox` varchar(50) DEFAULT '',
  `floor1` varchar(10) DEFAULT '',
  `room` varchar(10) DEFAULT '',
  `street_name` varchar(100) DEFAULT '',
  `additional_street_name` varchar(100) DEFAULT '',
  `block_name` varchar(100) DEFAULT '',
  `building_name` varchar(100) DEFAULT '',
  `building_number` varchar(100) DEFAULT '',
  `inhouse_mail` varchar(100) DEFAULT '',
  `department` varchar(100) DEFAULT '',
  `mark_attention` varchar(100) DEFAULT '',
  `mark_care` varchar(100) DEFAULT '',
  `plot_identification` varchar(50) DEFAULT '',
  `city_subdivision_name` varchar(100) DEFAULT '',
  `city_name` varchar(100) DEFAULT '',
  `postal_zone` varchar(20) DEFAULT '',
  `country_subentity` varchar(100) DEFAULT '',
  `country_subentity_code` varchar(20) DEFAULT '',
  `region` varchar(50) DEFAULT '',
  `district` varchar(50) DEFAULT '',
  `timezone_offset` varchar(50) DEFAULT '',
  `country_id_code` varchar(50) DEFAULT '',
  `country_name` varchar(100) DEFAULT '',
  `location_coord_lat` double DEFAULT 0,
  `location_coord_lon` double DEFAULT 0,
  `note` varchar(50) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `addresses`
--

LOCK TABLES `addresses` WRITE;
/*!40000 ALTER TABLE `addresses` DISABLE KEYS */;
/*!40000 ALTER TABLE `addresses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `allowance_charges`
--

DROP TABLE IF EXISTS `allowance_charges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `allowance_charges` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `ac_id` varchar(100) DEFAULT '',
  `charge_indicator` tinyint(1) DEFAULT 0,
  `allowance_charge_reason_code` varchar(100) DEFAULT '',
  `allowance_charge_reason` varchar(50) DEFAULT '',
  `multiplier_factor_numeric` int(10) unsigned DEFAULT 0,
  `prepaid_indicator` tinyint(1) DEFAULT 0,
  `sequence_numeric` int(10) unsigned DEFAULT 0,
  `amount` double DEFAULT 0,
  `base_amount` double DEFAULT 0,
  `per_unit_amount` double DEFAULT 0,
  `tax_category_id` int(10) unsigned DEFAULT 0,
  `tax_total_id` int(10) unsigned DEFAULT 0,
  `master_flag` varchar(20) DEFAULT '',
  `master_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `allowance_charges`
--

LOCK TABLES `allowance_charges` WRITE;
/*!40000 ALTER TABLE `allowance_charges` DISABLE KEYS */;
/*!40000 ALTER TABLE `allowance_charges` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bill_of_ladings`
--

DROP TABLE IF EXISTS `bill_of_ladings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bill_of_ladings` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `bill_of_lading_id` varchar(100) DEFAULT '',
  `carrier_assigned_id` varchar(100) DEFAULT '',
  `name1` varchar(100) DEFAULT '',
  `description` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `document_status_code` varchar(50) DEFAULT '',
  `shipping_order_id` varchar(100) DEFAULT '',
  `to_order_indicator` tinyint(1) DEFAULT 0,
  `ad_valorem_indicator` tinyint(1) DEFAULT 0,
  `declared_carriage_value_amount` double DEFAULT 0,
  `declared_carriage_value_amount_currency_code` varchar(20) DEFAULT '',
  `other_instruction` varchar(50) DEFAULT '',
  `consignor_party_id` int(10) unsigned DEFAULT 0,
  `carrier_party_id` int(10) unsigned DEFAULT 0,
  `freight_forwarder_party_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bill_of_ladings`
--

LOCK TABLES `bill_of_ladings` WRITE;
/*!40000 ALTER TABLE `bill_of_ladings` DISABLE KEYS */;
/*!40000 ALTER TABLE `bill_of_ladings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `consignments`
--

DROP TABLE IF EXISTS `consignments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `consignments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `cons_id` varchar(100) DEFAULT '',
  `carrier_assigned_id` varchar(100) DEFAULT '',
  `consignee_assigned_id` varchar(100) DEFAULT '',
  `consignor_assigned_id` varchar(100) DEFAULT '',
  `freight_forwarder_assigned_id` varchar(100) DEFAULT '',
  `broker_assigned_id` varchar(100) DEFAULT '',
  `contracted_carrier_assigned_id` varchar(100) DEFAULT '',
  `performing_carrier_assigned_id` varchar(100) DEFAULT '',
  `summary_description` varchar(50) DEFAULT '',
  `total_invoice_amount` double DEFAULT 0,
  `declared_customs_value_amount` double DEFAULT 0,
  `tariff_description` varchar(50) DEFAULT '',
  `tariff_code` varchar(50) DEFAULT '',
  `insurance_premium_amount` double DEFAULT 0,
  `gross_weight_measure` double DEFAULT 0,
  `net_weight_measure` double DEFAULT 0,
  `net_net_weight_measure` double DEFAULT 0,
  `chargeable_weight_measure` double DEFAULT 0,
  `gross_volume_measure` double DEFAULT 0,
  `net_volume_measure` double DEFAULT 0,
  `loading_length_measure` double DEFAULT 0,
  `remarks` varchar(50) DEFAULT '',
  `hazardous_risk_indicator` tinyint(1) DEFAULT 0,
  `animal_food_indicator` tinyint(1) DEFAULT 0,
  `human_food_indicator` tinyint(1) DEFAULT 0,
  `livestock_indicator` tinyint(1) DEFAULT 0,
  `bulk_cargo_indicator` tinyint(1) DEFAULT 0,
  `containerized_indicator` tinyint(1) DEFAULT 0,
  `general_cargo_indicator` tinyint(1) DEFAULT 0,
  `special_security_indicator` tinyint(1) DEFAULT 0,
  `third_party_payer_indicator` tinyint(1) DEFAULT 0,
  `carrier_service_instructions` varchar(50) DEFAULT '',
  `customs_clearance_service_instructions` varchar(50) DEFAULT '',
  `forwarder_service_instructions` varchar(50) DEFAULT '',
  `special_service_instructions` varchar(50) DEFAULT '',
  `sequence_id` int(10) unsigned DEFAULT 0,
  `shipping_priority_level_code` varchar(50) DEFAULT '',
  `handling_code` varchar(50) DEFAULT '',
  `handling_instructions` varchar(50) DEFAULT '',
  `information` varchar(50) DEFAULT '',
  `total_goods_item_quantity` int(10) unsigned DEFAULT 0,
  `total_transport_handling_unit_quantity` int(10) unsigned DEFAULT 0,
  `insurance_value_amount` double DEFAULT 0,
  `declared_for_carriage_value_amount` double DEFAULT 0,
  `declared_statistics_value_amount` double DEFAULT 0,
  `free_on_board_value_amount` double DEFAULT 0,
  `special_instructions` varchar(50) DEFAULT '',
  `split_consignment_indicator` tinyint(1) DEFAULT 0,
  `delivery_instructions` varchar(50) DEFAULT '',
  `consignment_quantity` double DEFAULT 0,
  `consolidatable_indicator` tinyint(1) DEFAULT 0,
  `haulage_instructions` varchar(255) DEFAULT NULL,
  `loading_sequence_id` int(10) unsigned DEFAULT 0,
  `child_consignment_quantity` int(10) unsigned DEFAULT 0,
  `total_packages_quantity` int(10) unsigned DEFAULT 0,
  `consignee_party_id` int(10) unsigned DEFAULT 0,
  `exporter_party_id` int(10) unsigned DEFAULT 0,
  `consignor_party_id` int(10) unsigned DEFAULT 0,
  `importer_party_id` int(10) unsigned DEFAULT 0,
  `carrier_party_id` int(10) unsigned DEFAULT 0,
  `freight_forwarder_party_id` int(10) unsigned DEFAULT 0,
  `notify_party_id` int(10) unsigned DEFAULT 0,
  `original_despatch_party_id` int(10) unsigned DEFAULT 0,
  `final_delivery_party_id` int(10) unsigned DEFAULT 0,
  `performing_carrier_party_id` int(10) unsigned DEFAULT 0,
  `substitute_carrier_party_id` int(10) unsigned DEFAULT 0,
  `logistics_operator_party_id` int(10) unsigned DEFAULT 0,
  `transport_advisor_party_id` int(10) unsigned DEFAULT 0,
  `hazardous_item_notification_party_id` int(10) unsigned DEFAULT 0,
  `insurance_party_id` int(10) unsigned DEFAULT 0,
  `mortgage_holder_party_id` int(10) unsigned DEFAULT 0,
  `bill_of_lading_holder_party_id` int(10) unsigned DEFAULT 0,
  `original_departure_country_id_code` varchar(50) DEFAULT '',
  `original_departure_country_name` varchar(50) DEFAULT '',
  `final_destination_country_id_code` varchar(50) DEFAULT '',
  `final_destination_country_name` varchar(50) DEFAULT '',
  `transit_country_id_code` varchar(50) DEFAULT '',
  `transit_country_name` varchar(50) DEFAULT '',
  `delivery_terms_id` int(10) unsigned DEFAULT 0,
  `payment_terms_id` int(10) unsigned DEFAULT 0,
  `collect_payment_terms_id` int(10) unsigned DEFAULT 0,
  `disbursement_payment_terms_id` int(10) unsigned DEFAULT 0,
  `prepaid_payment_terms_id` int(10) unsigned DEFAULT 0,
  `first_arrival_port_address_id` int(10) unsigned DEFAULT 0,
  `last_exit_port_location_address_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `consignments`
--

LOCK TABLES `consignments` WRITE;
/*!40000 ALTER TABLE `consignments` DISABLE KEYS */;
/*!40000 ALTER TABLE `consignments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `credit_note_headers`
--

DROP TABLE IF EXISTS `credit_note_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `credit_note_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `cnh_id` varchar(100) DEFAULT '',
  `credit_note_type_code` varchar(10) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `document_currency_code` varchar(20) DEFAULT '',
  `tax_currency_code` varchar(20) DEFAULT '',
  `pricing_currency_code` varchar(20) DEFAULT '',
  `payment_currency_code` varchar(20) DEFAULT '',
  `payment_alt_currency_code` varchar(20) DEFAULT '',
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `discrepancy_response` varchar(50) DEFAULT '',
  `order_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `despatch_id` int(10) unsigned DEFAULT 0,
  `receipt_id` int(10) unsigned DEFAULT 0,
  `contract_id` int(10) unsigned DEFAULT 0,
  `statement_id` int(10) unsigned DEFAULT 0,
  `signature` varchar(100) DEFAULT '',
  `accounting_supplier_party_id` int(10) unsigned DEFAULT 0,
  `accounting_customer_party_id` int(10) unsigned DEFAULT 0,
  `payee_party_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `tax_representative_party_id` int(10) unsigned DEFAULT 0,
  `tax_ex_source_currency_code` varchar(20) DEFAULT '',
  `tax_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_target_currency_code` varchar(20) DEFAULT '',
  `tax_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `tax_ex_calculation_rate` double DEFAULT 0,
  `tax_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `pricing_ex_calculation_rate` double DEFAULT 0,
  `pricing_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_ex_calculation_rate` double DEFAULT 0,
  `payment_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_alt_ex_calculation_rate` double DEFAULT 0,
  `payment_alt_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `line_extension_amount` double DEFAULT 0,
  `tax_exclusive_amount` double DEFAULT 0,
  `tax_inclusive_amount` double DEFAULT 0,
  `allowance_total_amount` double DEFAULT 0,
  `charge_total_amount` double DEFAULT 0,
  `withholding_tax_total_amount` double DEFAULT 0,
  `prepaid_amount` double DEFAULT 0,
  `payable_rounding_amount` double DEFAULT 0,
  `payable_amount` double DEFAULT 0,
  `payable_alternative_amount` double DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `due_date` datetime DEFAULT current_timestamp(),
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `invoice_period_start_date` datetime DEFAULT current_timestamp(),
  `invoice_period_end_date` datetime DEFAULT current_timestamp(),
  `tax_ex_date` datetime DEFAULT current_timestamp(),
  `pricing_ex_date` datetime DEFAULT current_timestamp(),
  `payment_ex_date` datetime DEFAULT current_timestamp(),
  `payment_alt_ex_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `credit_note_headers`
--

LOCK TABLES `credit_note_headers` WRITE;
/*!40000 ALTER TABLE `credit_note_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `credit_note_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `credit_note_lines`
--

DROP TABLE IF EXISTS `credit_note_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `credit_note_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `cnl_id` varchar(100) DEFAULT '',
  `note` varchar(100) DEFAULT '',
  `credited_quantity` double DEFAULT 0,
  `line_extension_amount` double DEFAULT 0,
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `payment_purpose_code` varchar(10) DEFAULT '',
  `free_of_charge_indicator` tinyint(1) DEFAULT 0,
  `discrepancy_response` varchar(50) DEFAULT '',
  `order_line_id` int(10) unsigned DEFAULT 0,
  `despatch_line_id` int(10) unsigned DEFAULT 0,
  `receipt_line_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `originator_party_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `price_amount` double DEFAULT 0,
  `price_base_quantity` double DEFAULT 0,
  `price_change_reason` varchar(100) DEFAULT '',
  `price_type_code` varchar(20) DEFAULT '',
  `price_type` varchar(20) DEFAULT '',
  `orderable_unit_factor_rate` double DEFAULT 0,
  `price_list_id` int(10) unsigned DEFAULT 0,
  `credit_note_header_id` int(10) unsigned DEFAULT 0,
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `invoice_period_start_date` datetime DEFAULT current_timestamp(),
  `invoice_period_end_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_start_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `credit_note_lines`
--

LOCK TABLES `credit_note_lines` WRITE;
/*!40000 ALTER TABLE `credit_note_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `credit_note_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `debit_note_headers`
--

DROP TABLE IF EXISTS `debit_note_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `debit_note_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `dnh_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `document_currency_code` varchar(20) DEFAULT '',
  `tax_currency_code` varchar(20) DEFAULT '',
  `pricing_currency_code` varchar(20) DEFAULT '',
  `payment_currency_code` varchar(20) DEFAULT '',
  `payment_alt_currency_code` varchar(20) DEFAULT '',
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `discrepancy_response` varchar(50) DEFAULT '',
  `order_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `despatch_id` int(10) unsigned DEFAULT 0,
  `receipt_id` int(10) unsigned DEFAULT 0,
  `statement_id` int(10) unsigned DEFAULT 0,
  `contract_id` int(10) unsigned DEFAULT 0,
  `accounting_supplier_party_id` int(10) unsigned DEFAULT 0,
  `accounting_customer_party_id` int(10) unsigned DEFAULT 0,
  `payee_party_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `tax_representative_party_id` int(10) unsigned DEFAULT 0,
  `tax_ex_source_currency_code` varchar(20) DEFAULT '',
  `tax_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_target_currency_code` varchar(20) DEFAULT '',
  `tax_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `tax_ex_calculation_rate` double DEFAULT 0,
  `tax_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `pricing_ex_calculation_rate` double DEFAULT 0,
  `pricing_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_ex_calculation_rate` double DEFAULT 0,
  `payment_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_alt_ex_calculation_rate` double DEFAULT 0,
  `payment_alt_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `line_extension_amount` double DEFAULT 0,
  `tax_exclusive_amount` double DEFAULT 0,
  `tax_inclusive_amount` double DEFAULT 0,
  `allowance_total_amount` double DEFAULT 0,
  `charge_total_amount` double DEFAULT 0,
  `withholding_tax_total_amount` double DEFAULT 0,
  `prepaid_amount` double DEFAULT 0,
  `payable_rounding_amount` double DEFAULT 0,
  `payable_amount` double DEFAULT 0,
  `payable_alternative_amount` double DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `invoice_period_start_date` datetime DEFAULT current_timestamp(),
  `invoice_period_end_date` datetime DEFAULT current_timestamp(),
  `tax_ex_date` datetime DEFAULT current_timestamp(),
  `pricing_ex_date` datetime DEFAULT current_timestamp(),
  `payment_ex_date` datetime DEFAULT current_timestamp(),
  `payment_alt_ex_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `debit_note_headers`
--

LOCK TABLES `debit_note_headers` WRITE;
/*!40000 ALTER TABLE `debit_note_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `debit_note_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `debit_note_lines`
--

DROP TABLE IF EXISTS `debit_note_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `debit_note_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `dnl_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `debited_quantity` double DEFAULT 0,
  `line_extension_amount` double DEFAULT 0,
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `payment_purpose_code` varchar(20) DEFAULT '',
  `discrepancy_response` varchar(50) DEFAULT '',
  `despatch_line_id` int(10) unsigned DEFAULT 0,
  `receipt_line_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `price_amount` double DEFAULT 0,
  `price_base_quantity` double DEFAULT 0,
  `price_change_reason` varchar(100) DEFAULT '',
  `price_type_code` varchar(20) DEFAULT '',
  `price_type` varchar(20) DEFAULT '',
  `orderable_unit_factor_rate` double DEFAULT 0,
  `price_list_id` int(10) unsigned DEFAULT 0,
  `debit_note_header_id` int(10) unsigned DEFAULT 0,
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_start_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `debit_note_lines`
--

LOCK TABLES `debit_note_lines` WRITE;
/*!40000 ALTER TABLE `debit_note_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `debit_note_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deliveries`
--

DROP TABLE IF EXISTS `deliveries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deliveries` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `del_id` varchar(100) DEFAULT '',
  `quantity` double DEFAULT 0,
  `minimum_quantity` double DEFAULT 0,
  `maximum_quantity` double DEFAULT 0,
  `release_id` varchar(100) DEFAULT '',
  `tracking_id` varchar(100) DEFAULT '',
  `minimum_batch_quantity` bigint(20) DEFAULT 0,
  `maximum_batch_quantity` bigint(20) DEFAULT 0,
  `consumer_unit_quantity` bigint(20) DEFAULT 0,
  `hazardous_risk_indicator` tinyint(1) DEFAULT 0,
  `delivery_address_id` int(10) unsigned DEFAULT 0,
  `delivery_location_id` int(10) unsigned DEFAULT 0,
  `alternative_delivery_location_id` int(10) unsigned DEFAULT 0,
  `carrier_party_id` int(10) unsigned DEFAULT 0,
  `delivery_party_id` int(10) unsigned DEFAULT 0,
  `notify_party_id` int(10) unsigned DEFAULT 0,
  `despatch_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `actual_delivery_date` datetime DEFAULT current_timestamp(),
  `latest_delivery_date` datetime DEFAULT current_timestamp(),
  `requested_delivery_period_start_date` datetime DEFAULT current_timestamp(),
  `requested_delivery_period_end_date` datetime DEFAULT current_timestamp(),
  `promised_delivery_period_start_date` datetime DEFAULT current_timestamp(),
  `promised_delivery_period_end_date` datetime DEFAULT current_timestamp(),
  `estimated_delivery_period_start_date` datetime DEFAULT current_timestamp(),
  `estimated_delivery_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deliveries`
--

LOCK TABLES `deliveries` WRITE;
/*!40000 ALTER TABLE `deliveries` DISABLE KEYS */;
/*!40000 ALTER TABLE `deliveries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `delivery_terms`
--

DROP TABLE IF EXISTS `delivery_terms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `delivery_terms` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `del_term_id` varchar(100) DEFAULT '',
  `special_terms` varchar(50) DEFAULT '',
  `loss_risk_responsibility_code` varchar(100) DEFAULT '',
  `loss_risk` varchar(50) DEFAULT '',
  `amount` double DEFAULT 0,
  `delivery_location_id` int(10) unsigned DEFAULT 0,
  `del_term_allowance_charge_id` int(10) unsigned DEFAULT 0,
  `delivery_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `delivery_terms`
--

LOCK TABLES `delivery_terms` WRITE;
/*!40000 ALTER TABLE `delivery_terms` DISABLE KEYS */;
/*!40000 ALTER TABLE `delivery_terms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `despatch_headers`
--

DROP TABLE IF EXISTS `despatch_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `despatch_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `desph_id` varchar(100) DEFAULT '',
  `document_status_code` varchar(50) DEFAULT '',
  `despatch_advice_type_code` varchar(50) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `order_id` int(10) unsigned DEFAULT 0,
  `despatch_supplier_party_id` int(10) unsigned DEFAULT 0,
  `delivery_customer_party_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `originator_customer_party_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `despatch_headers`
--

LOCK TABLES `despatch_headers` WRITE;
/*!40000 ALTER TABLE `despatch_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `despatch_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `despatch_lines`
--

DROP TABLE IF EXISTS `despatch_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `despatch_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `despl_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `line_status_code` varchar(50) DEFAULT '',
  `delivered_quantity` double DEFAULT 0,
  `backorder_quantity` double DEFAULT 0,
  `backorder_reason` varchar(50) DEFAULT '',
  `outstanding_quantity` double DEFAULT 0,
  `outstanding_reason` varchar(50) DEFAULT '',
  `oversupply_quantity` double DEFAULT 0,
  `order_line_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `despatch_header_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `despatch_lines`
--

LOCK TABLES `despatch_lines` WRITE;
/*!40000 ALTER TABLE `despatch_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `despatch_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `despatches`
--

DROP TABLE IF EXISTS `despatches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `despatches` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `desp_id` varchar(100) DEFAULT '',
  `release_id` varchar(100) DEFAULT '',
  `instructions` varchar(50) DEFAULT '',
  `despatch_address_id` int(10) unsigned DEFAULT 0,
  `despatch_location_id` int(10) unsigned DEFAULT 0,
  `despatch_party_contact_id` int(10) unsigned DEFAULT 0,
  `despatch_party_id` int(10) unsigned DEFAULT 0,
  `carrier_party_id` int(10) unsigned DEFAULT 0,
  `notify_party_id` int(10) unsigned DEFAULT 0,
  `requested_despatch_date` datetime DEFAULT current_timestamp(),
  `estimated_despatch_date` datetime DEFAULT current_timestamp(),
  `actual_despatch_date` datetime DEFAULT current_timestamp(),
  `guaranteed_despatch_date` datetime DEFAULT current_timestamp(),
  `estimated_despatch_period_start_date` datetime DEFAULT current_timestamp(),
  `estimated_despatch_period_end_date` datetime DEFAULT current_timestamp(),
  `requested_despatch_period_start_date` datetime DEFAULT current_timestamp(),
  `requested_despatch_period_end_date` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `despatches`
--

LOCK TABLES `despatches` WRITE;
/*!40000 ALTER TABLE `despatches` DISABLE KEYS */;
/*!40000 ALTER TABLE `despatches` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `financial_institution_branches`
--

DROP TABLE IF EXISTS `financial_institution_branches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `financial_institution_branches` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `fb_id` varchar(100) DEFAULT '',
  `name1` varchar(200) DEFAULT NULL,
  `address_id` int(10) unsigned DEFAULT 0,
  `financial_institution_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `financial_institution_branches`
--

LOCK TABLES `financial_institution_branches` WRITE;
/*!40000 ALTER TABLE `financial_institution_branches` DISABLE KEYS */;
/*!40000 ALTER TABLE `financial_institution_branches` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `financial_institutions`
--

DROP TABLE IF EXISTS `financial_institutions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `financial_institutions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `fi_id` varchar(100) DEFAULT '',
  `name1` varchar(200) DEFAULT NULL,
  `address_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `financial_institutions`
--

LOCK TABLES `financial_institutions` WRITE;
/*!40000 ALTER TABLE `financial_institutions` DISABLE KEYS */;
/*!40000 ALTER TABLE `financial_institutions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoice_headers`
--

DROP TABLE IF EXISTS `invoice_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoice_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `ih_id` varchar(100) DEFAULT '',
  `invoice_type_code` varchar(50) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `document_currency_code` varchar(20) DEFAULT '',
  `tax_currency_code` varchar(20) DEFAULT '',
  `pricing_currency_code` varchar(20) DEFAULT '',
  `payment_currency_code` varchar(20) DEFAULT '',
  `payment_alt_currency_code` varchar(20) DEFAULT '',
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `order_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `despatch_id` int(10) unsigned DEFAULT 0,
  `receipt_id` int(10) unsigned DEFAULT 0,
  `statement_id` int(10) unsigned DEFAULT 0,
  `contract_id` int(10) unsigned DEFAULT 0,
  `accounting_supplier_party_id` int(10) unsigned DEFAULT 0,
  `accounting_customer_party_id` int(10) unsigned DEFAULT 0,
  `payee_party_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `tax_representative_party_id` int(10) unsigned DEFAULT 0,
  `tax_ex_source_currency_code` varchar(20) DEFAULT '',
  `tax_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_target_currency_code` varchar(20) DEFAULT '',
  `tax_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `tax_ex_calculation_rate` double DEFAULT 0,
  `tax_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `pricing_ex_calculation_rate` double DEFAULT 0,
  `pricing_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_ex_calculation_rate` double DEFAULT 0,
  `payment_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_alt_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_alt_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_alt_ex_calculation_rate` double DEFAULT 0,
  `payment_alt_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `line_extension_amount` double DEFAULT 0,
  `tax_exclusive_amount` double DEFAULT 0,
  `tax_inclusive_amount` double DEFAULT 0,
  `allowance_total_amount` double DEFAULT 0,
  `charge_total_amount` double DEFAULT 0,
  `withholding_tax_total_amount` double DEFAULT 0,
  `prepaid_amount` double DEFAULT 0,
  `payable_rounding_amount` double DEFAULT 0,
  `payable_amount` double DEFAULT 0,
  `payable_alternative_amount` double DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `due_date` datetime DEFAULT current_timestamp(),
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `invoice_period_start_date` datetime DEFAULT current_timestamp(),
  `invoice_period_end_date` datetime DEFAULT current_timestamp(),
  `tax_ex_date` datetime DEFAULT current_timestamp(),
  `pricing_ex_date` datetime DEFAULT current_timestamp(),
  `payment_ex_date` datetime DEFAULT current_timestamp(),
  `payment_alt_ex_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoice_headers`
--

LOCK TABLES `invoice_headers` WRITE;
/*!40000 ALTER TABLE `invoice_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `invoice_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoice_lines`
--

DROP TABLE IF EXISTS `invoice_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invoice_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `il_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `invoiced_quantity` double DEFAULT 0,
  `line_extension_amount` double DEFAULT 0,
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `payment_purpose_code` varchar(50) DEFAULT '',
  `free_of_charge_indicator` tinyint(1) DEFAULT 0,
  `order_line_id` int(10) unsigned DEFAULT 0,
  `despatch_line_id` int(10) unsigned DEFAULT 0,
  `receipt_line_id` int(10) unsigned DEFAULT 0,
  `billing_id` int(10) unsigned DEFAULT 0,
  `originator_party_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `price_amount` double DEFAULT 0,
  `price_base_quantity` double DEFAULT 0,
  `price_change_reason` varchar(100) DEFAULT '',
  `price_type_code` varchar(20) DEFAULT '',
  `price_type` varchar(20) DEFAULT '',
  `orderable_unit_factor_rate` double DEFAULT 0,
  `price_list_id` int(10) unsigned DEFAULT 0,
  `invoice_header_id` int(10) unsigned DEFAULT 0,
  `tax_point_date` datetime DEFAULT current_timestamp(),
  `invoice_period_start_date` datetime DEFAULT current_timestamp(),
  `invoice_period_end_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_start_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoice_lines`
--

LOCK TABLES `invoice_lines` WRITE;
/*!40000 ALTER TABLE `invoice_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `invoice_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_certificates`
--

DROP TABLE IF EXISTS `item_certificates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_certificates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cert_id` varchar(100) DEFAULT '',
  `certificate_type_code` varchar(50) DEFAULT '',
  `certificate_type` varchar(200) DEFAULT NULL,
  `remarks` varchar(50) DEFAULT '',
  `party_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_certificates`
--

LOCK TABLES `item_certificates` WRITE;
/*!40000 ALTER TABLE `item_certificates` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_certificates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_commodity_classifications`
--

DROP TABLE IF EXISTS `item_commodity_classifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_commodity_classifications` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `nature_code` varchar(100) DEFAULT '',
  `cargo_type_code` varchar(100) DEFAULT '',
  `commodity_code` varchar(100) DEFAULT '',
  `item_classification_code` varchar(100) DEFAULT '',
  `item_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_commodity_classifications`
--

LOCK TABLES `item_commodity_classifications` WRITE;
/*!40000 ALTER TABLE `item_commodity_classifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_commodity_classifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_dimensions`
--

DROP TABLE IF EXISTS `item_dimensions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_dimensions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `attribute_id` varchar(100) DEFAULT '',
  `measure` double DEFAULT 0,
  `description` varchar(50) DEFAULT '',
  `minimum_measure` double DEFAULT 0,
  `maximum_measure` double DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_dimensions`
--

LOCK TABLES `item_dimensions` WRITE;
/*!40000 ALTER TABLE `item_dimensions` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_dimensions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_instances`
--

DROP TABLE IF EXISTS `item_instances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_instances` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `product_trace_id` varchar(100) DEFAULT '',
  `registration_id` varchar(100) DEFAULT '',
  `serial_id` varchar(100) DEFAULT '',
  `lot_number_id` varchar(100) DEFAULT '',
  `item_id` int(10) unsigned DEFAULT 0,
  `manufacture_date` datetime DEFAULT current_timestamp(),
  `best_before_date` datetime DEFAULT current_timestamp(),
  `lot_expiry_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_instances`
--

LOCK TABLES `item_instances` WRITE;
/*!40000 ALTER TABLE `item_instances` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_instances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_properties`
--

DROP TABLE IF EXISTS `item_properties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_properties` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `item_property_id` varchar(100) DEFAULT '',
  `item_property_name` varchar(200) DEFAULT NULL,
  `item_property_name_code` varchar(50) DEFAULT '',
  `test_method` varchar(100) DEFAULT '',
  `value` varchar(200) DEFAULT NULL,
  `value_quantity` double DEFAULT 0,
  `value_qualifier` varchar(100) DEFAULT '',
  `importance_code` varchar(50) DEFAULT '',
  `list_value` varchar(50) DEFAULT '',
  `item_property_range_measure` double DEFAULT 0,
  `item_property_range_min_value` double DEFAULT 0,
  `item_property_range_max_value` double DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `usability_period_start_date` datetime DEFAULT current_timestamp(),
  `usability_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_properties`
--

LOCK TABLES `item_properties` WRITE;
/*!40000 ALTER TABLE `item_properties` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_properties` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `item_property_groups`
--

DROP TABLE IF EXISTS `item_property_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `item_property_groups` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `item_property_group_id` varchar(100) DEFAULT '',
  `item_property_group_name` varchar(200) DEFAULT NULL,
  `item_property_group_importance_code` varchar(50) DEFAULT '',
  `item_property_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `item_property_groups`
--

LOCK TABLES `item_property_groups` WRITE;
/*!40000 ALTER TABLE `item_property_groups` DISABLE KEYS */;
/*!40000 ALTER TABLE `item_property_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `items`
--

DROP TABLE IF EXISTS `items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `items` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `description` varchar(50) DEFAULT '',
  `pack_quantity` bigint(20) DEFAULT 0,
  `pack_size_numeric` bigint(20) DEFAULT 0,
  `catalogue_indicator` tinyint(1) DEFAULT 0,
  `item_name` varchar(50) DEFAULT '',
  `hazardous_risk_indicator` tinyint(1) DEFAULT 0,
  `additional_information` varchar(50) DEFAULT '',
  `keyword` varchar(50) DEFAULT '',
  `brand_name` varchar(50) DEFAULT '',
  `model_name` varchar(50) DEFAULT '',
  `buyers_item_identification_id` varchar(100) DEFAULT '',
  `sellers_item_identification_id` varchar(100) DEFAULT '',
  `manufacturers_item_identification_id` varchar(100) DEFAULT '',
  `standard_item_identification_id` varchar(100) DEFAULT '',
  `catalogue_item_identification_id` varchar(100) DEFAULT '',
  `additional_item_identification_id` varchar(100) DEFAULT '',
  `origin_country_id_code` varchar(50) DEFAULT '',
  `origin_country_name` varchar(100) DEFAULT '',
  `manufacturer_party_id` int(10) unsigned DEFAULT 0,
  `information_content_provider_party_id` int(10) unsigned DEFAULT 0,
  `tax_category_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `items`
--

LOCK TABLES `items` WRITE;
/*!40000 ALTER TABLE `items` DISABLE KEYS */;
/*!40000 ALTER TABLE `items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `locations`
--

DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `loc_id` varchar(100) DEFAULT '',
  `description` varchar(100) DEFAULT '',
  `conditions` varchar(100) DEFAULT '',
  `country_subentity` varchar(100) DEFAULT '',
  `country_subentity_code` varchar(100) DEFAULT '',
  `location_type_code` varchar(100) DEFAULT '',
  `information_uri` varchar(100) DEFAULT '',
  `loc_name` varchar(100) DEFAULT '',
  `location_coord_lat` varchar(10) DEFAULT '',
  `location_coord_lon` varchar(11) DEFAULT '',
  `altitude_measure` double DEFAULT 0,
  `address_id` int(10) unsigned DEFAULT 0,
  `validity_period_start_date` datetime DEFAULT current_timestamp(),
  `validity_period_end_date` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `locations`
--

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `parties`
--

DROP TABLE IF EXISTS `parties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `parties` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `party_endpoint_id` varchar(50) DEFAULT '',
  `party_endpoint_scheme_id` varchar(50) DEFAULT '',
  `party_name` varchar(200) DEFAULT '',
  `party_desc` varchar(200) DEFAULT '',
  `party_type` varchar(50) DEFAULT '',
  `registration_name` varchar(200) DEFAULT '',
  `company_id` varchar(50) DEFAULT '',
  `company_legal_form_code` varchar(20) DEFAULT '',
  `company_legal_form` varchar(50) DEFAULT '',
  `sole_proprietorship_indicator` tinyint(1) DEFAULT 0,
  `company_liquidation_status_code` varchar(20) DEFAULT '',
  `corporate_stock_amount` bigint(20) DEFAULT 0,
  `fully_paid_shares_indicator` tinyint(1) DEFAULT 0,
  `corporate_registration_id` varchar(50) DEFAULT '',
  `corporate_registration_name` varchar(200) DEFAULT '',
  `corporate_registration_type_code` varchar(20) DEFAULT '',
  `tax_level_code` varchar(20) DEFAULT '',
  `exemption_reason_code` varchar(20) DEFAULT '',
  `exemption_reason` varchar(50) DEFAULT '',
  `tax_scheme_id` int(10) unsigned DEFAULT 0,
  `registration_date` datetime DEFAULT current_timestamp(),
  `registration_expiration_date` datetime DEFAULT current_timestamp(),
  `level_p` int(10) unsigned DEFAULT 0,
  `parent_id` int(10) unsigned DEFAULT 0,
  `num_chd` smallint(6) DEFAULT 0,
  `leaf` tinyint(1) DEFAULT 0,
  `tax_reference1` varchar(50) DEFAULT '',
  `tax_reference2` varchar(50) DEFAULT '',
  `public_key` varchar(500) DEFAULT '',
  `address_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `parties`
--

LOCK TABLES `parties` WRITE;
/*!40000 ALTER TABLE `parties` DISABLE KEYS */;
/*!40000 ALTER TABLE `parties` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_chds`
--

DROP TABLE IF EXISTS `party_chds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_chds` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `party_id` int(10) unsigned DEFAULT 0,
  `party_chd_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_chds`
--

LOCK TABLES `party_chds` WRITE;
/*!40000 ALTER TABLE `party_chds` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_chds` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_contact_rels`
--

DROP TABLE IF EXISTS `party_contact_rels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_contact_rels` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `party_id` int(10) unsigned DEFAULT 0,
  `party_contact_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_contact_rels`
--

LOCK TABLES `party_contact_rels` WRITE;
/*!40000 ALTER TABLE `party_contact_rels` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_contact_rels` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_contacts`
--

DROP TABLE IF EXISTS `party_contacts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_contacts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `first_name` varchar(255) NOT NULL,
  `middle_name` varchar(50) DEFAULT '',
  `last_name` varchar(255) NOT NULL,
  `title` varchar(50) DEFAULT '',
  `name_suffix` varchar(10) DEFAULT '',
  `job_title` varchar(100) DEFAULT '',
  `org_dept` varchar(100) DEFAULT '',
  `email` varchar(100) DEFAULT '',
  `phone_mobile` varchar(20) DEFAULT '',
  `phone_work` varchar(20) DEFAULT '',
  `phone_fax` varchar(20) DEFAULT '',
  `country_calling_code` varchar(255) DEFAULT '+91',
  `url` varchar(500) DEFAULT '',
  `gender_code` varchar(10) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `user_id` int(10) unsigned DEFAULT 0,
  `party_id` int(10) unsigned DEFAULT 0,
  `party_name` varchar(200) DEFAULT NULL,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_contacts`
--

LOCK TABLES `party_contacts` WRITE;
/*!40000 ALTER TABLE `party_contacts` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_contacts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_corporate_jurisdictions`
--

DROP TABLE IF EXISTS `party_corporate_jurisdictions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_corporate_jurisdictions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `address_id` int(10) unsigned DEFAULT 0,
  `party_id` int(10) unsigned DEFAULT 0,
  `party_name` varchar(200) DEFAULT NULL,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_corporate_jurisdictions`
--

LOCK TABLES `party_corporate_jurisdictions` WRITE;
/*!40000 ALTER TABLE `party_corporate_jurisdictions` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_corporate_jurisdictions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_identifications`
--

DROP TABLE IF EXISTS `party_identifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_identifications` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `party_identification` varchar(200) DEFAULT '',
  `party_identification_scheme_id` varchar(200) DEFAULT '',
  `party_identification_scheme_name` varchar(200) DEFAULT '',
  `party_id` int(10) unsigned DEFAULT 0,
  `party_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_identifications`
--

LOCK TABLES `party_identifications` WRITE;
/*!40000 ALTER TABLE `party_identifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_identifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `party_social_profiles`
--

DROP TABLE IF EXISTS `party_social_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `party_social_profiles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `social_profle_name` varchar(50) DEFAULT '',
  `social_media_type_code` varchar(10) DEFAULT '',
  `uri` varchar(500) DEFAULT '',
  `party_id` int(10) unsigned DEFAULT 0,
  `party_name` varchar(200) DEFAULT NULL,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `party_social_profiles`
--

LOCK TABLES `party_social_profiles` WRITE;
/*!40000 ALTER TABLE `party_social_profiles` DISABLE KEYS */;
/*!40000 ALTER TABLE `party_social_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_mandate_clause_contents`
--

DROP TABLE IF EXISTS `payment_mandate_clause_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_mandate_clause_contents` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) DEFAULT NULL,
  `payment_mandate_clause_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_mandate_clause_contents`
--

LOCK TABLES `payment_mandate_clause_contents` WRITE;
/*!40000 ALTER TABLE `payment_mandate_clause_contents` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_mandate_clause_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_mandate_clauses`
--

DROP TABLE IF EXISTS `payment_mandate_clauses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_mandate_clauses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pm_cl_id` varchar(100) DEFAULT '',
  `payment_mandate_id` int(10) unsigned DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_mandate_clauses`
--

LOCK TABLES `payment_mandate_clauses` WRITE;
/*!40000 ALTER TABLE `payment_mandate_clauses` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_mandate_clauses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_mandates`
--

DROP TABLE IF EXISTS `payment_mandates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_mandates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `pmd_id` varchar(100) DEFAULT '',
  `mandate_type_code` varchar(50) DEFAULT '',
  `maximum_payment_instructions_numeric` int(10) unsigned DEFAULT 0,
  `maximum_paid_amount` double DEFAULT 0,
  `signature_id` varchar(100) DEFAULT '',
  `payer_party_id` int(10) unsigned DEFAULT 0,
  `payer_financial_account_id` int(10) unsigned DEFAULT 0,
  `clause` varchar(50) DEFAULT '',
  `validity_period_start_date` datetime DEFAULT current_timestamp(),
  `validity_period_end_date` datetime DEFAULT current_timestamp(),
  `payment_reversal_period_start_date` datetime DEFAULT current_timestamp(),
  `payment_reversal_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_mandates`
--

LOCK TABLES `payment_mandates` WRITE;
/*!40000 ALTER TABLE `payment_mandates` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_mandates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_means`
--

DROP TABLE IF EXISTS `payment_means`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_means` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `pm_id` varchar(100) DEFAULT '',
  `payment_means_code` varchar(50) DEFAULT '',
  `payment_channel_code` varchar(50) DEFAULT '',
  `instruction_id` varchar(100) DEFAULT '',
  `instruction_note` varchar(50) DEFAULT '',
  `credit_account_id` int(10) unsigned DEFAULT 0,
  `payment_term_id` int(10) unsigned DEFAULT 0,
  `payment_mandate_id` int(10) unsigned DEFAULT 0,
  `trade_financing_id` int(10) unsigned DEFAULT 0,
  `payer_financial_account_id` int(10) unsigned DEFAULT 0,
  `payee_financial_account_id` int(10) unsigned DEFAULT 0,
  `payment_due_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_means`
--

LOCK TABLES `payment_means` WRITE;
/*!40000 ALTER TABLE `payment_means` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_means` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_terms`
--

DROP TABLE IF EXISTS `payment_terms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payment_terms` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `pt_id` varchar(100) DEFAULT '',
  `prepaid_payment_reference_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `reference_event_code` varchar(50) DEFAULT '',
  `settlement_discount_percent` double DEFAULT 0,
  `penalty_surcharge_percent` double DEFAULT 0,
  `payment_percent` double DEFAULT 0,
  `amount` double DEFAULT 0,
  `settlement_discount_amount` double DEFAULT 0,
  `penalty_amount` double DEFAULT 0,
  `payment_terms_details_uri` varchar(50) DEFAULT '',
  `payment_means_id` int(10) unsigned DEFAULT 0,
  `payment_due_date` datetime DEFAULT current_timestamp(),
  `installment_due_date` datetime DEFAULT current_timestamp(),
  `settlement_period_start_date` datetime DEFAULT current_timestamp(),
  `settlement_period_end_date` datetime DEFAULT current_timestamp(),
  `penalty_period_start_date` datetime DEFAULT current_timestamp(),
  `penalty_period_end_date` datetime DEFAULT current_timestamp(),
  `exchange_rate` double DEFAULT 0,
  `validity_period_start_date` datetime DEFAULT current_timestamp(),
  `validity_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_terms`
--

LOCK TABLES `payment_terms` WRITE;
/*!40000 ALTER TABLE `payment_terms` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_terms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payments`
--

DROP TABLE IF EXISTS `payments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `payments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `p_id` varchar(100) DEFAULT '',
  `paid_amount` double DEFAULT 0,
  `instruction_id` varchar(100) DEFAULT '',
  `payment_mean_id` int(10) unsigned DEFAULT 0,
  `received_date` datetime DEFAULT current_timestamp(),
  `paid_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payments`
--

LOCK TABLES `payments` WRITE;
/*!40000 ALTER TABLE `payments` DISABLE KEYS */;
/*!40000 ALTER TABLE `payments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `purchase_order_headers`
--

DROP TABLE IF EXISTS `purchase_order_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_order_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `poh_id` varchar(100) DEFAULT '',
  `sales_order_id` varchar(100) DEFAULT '',
  `order_type_code` varchar(50) DEFAULT '',
  `note` varchar(255) DEFAULT '',
  `requested_invoice_currency_code` varchar(20) DEFAULT '',
  `document_currency_code` varchar(20) DEFAULT '',
  `pricing_currency_code` varchar(20) DEFAULT '',
  `tax_currency_code` varchar(20) DEFAULT '',
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `quotation_id` int(10) unsigned DEFAULT 0,
  `order_id` int(10) unsigned DEFAULT 0,
  `catalogue_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `originator_customer_party_id` int(10) unsigned DEFAULT 0,
  `freight_forwarder_party_id` int(10) unsigned DEFAULT 0,
  `accounting_customer_party_id` int(10) unsigned DEFAULT 0,
  `transaction_conditions` varchar(50) DEFAULT '',
  `tax_ex_source_currency_code` varchar(20) DEFAULT '',
  `tax_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_target_currency_code` varchar(20) DEFAULT '',
  `tax_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `tax_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `tax_ex_calculation_rate` double DEFAULT 0,
  `tax_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_code` varchar(20) DEFAULT '',
  `pricing_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `pricing_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `pricing_ex_calculation_rate` double DEFAULT 0,
  `pricing_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_code` varchar(20) DEFAULT '',
  `payment_ex_source_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_target_currency_code` varchar(20) DEFAULT '',
  `payment_ex_target_currency_base_rate` varchar(20) DEFAULT '',
  `payment_ex_exchange_market_id` int(10) unsigned DEFAULT 0,
  `payment_ex_calculation_rate` double DEFAULT 0,
  `payment_ex_mathematic_operator_code` varchar(20) DEFAULT '',
  `destination_country` varchar(50) DEFAULT '',
  `line_extension_amount` double DEFAULT 0,
  `tax_exclusive_amount` double DEFAULT 0,
  `tax_inclusive_amount` double DEFAULT 0,
  `allowance_total_amount` double DEFAULT 0,
  `charge_total_amount` double DEFAULT 0,
  `withholding_tax_total_amount` double DEFAULT 0,
  `prepaid_amount` double DEFAULT 0,
  `payable_rounding_amount` double DEFAULT 0,
  `payable_amount` double DEFAULT 0,
  `payable_alternative_amount` double DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `validity_period` datetime DEFAULT current_timestamp(),
  `tax_ex_date` datetime DEFAULT current_timestamp(),
  `pricing_ex_date` datetime DEFAULT current_timestamp(),
  `payment_ex_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `purchase_order_headers`
--

LOCK TABLES `purchase_order_headers` WRITE;
/*!40000 ALTER TABLE `purchase_order_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `purchase_order_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `purchase_order_lines`
--

DROP TABLE IF EXISTS `purchase_order_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `purchase_order_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `pol_id` varchar(100) DEFAULT '',
  `substitution_status_code` varchar(10) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `sales_order_id` varchar(100) DEFAULT '',
  `line_status_code` varchar(50) DEFAULT '',
  `quantity` double DEFAULT 0,
  `line_extension_amount` double DEFAULT 0,
  `total_tax_amount` double DEFAULT 0,
  `minimum_quantity` double DEFAULT 0,
  `maximum_quantity` double DEFAULT 0,
  `minimum_backorder_quantity` double DEFAULT 0,
  `maximum_backorder_quantity` double DEFAULT 0,
  `inspection_method_code` varchar(50) DEFAULT '',
  `partial_delivery_indicator` tinyint(1) DEFAULT 0,
  `back_order_allowed_indicator` tinyint(1) DEFAULT 0,
  `accounting_cost_code` varchar(50) DEFAULT '',
  `accounting_cost` varchar(100) DEFAULT '',
  `warranty_information` varchar(50) DEFAULT '',
  `originator_party_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `price_amount` double DEFAULT 0,
  `price_base_quantity` double DEFAULT 0,
  `price_change_reason` varchar(100) DEFAULT '',
  `price_type_code` varchar(20) DEFAULT '',
  `price_type` varchar(20) DEFAULT '',
  `orderable_unit_factor_rate` double DEFAULT 0,
  `price_list_id` int(10) unsigned DEFAULT 0,
  `purchase_order_header_id` int(10) unsigned DEFAULT 0,
  `price_validity_period_start_date` datetime DEFAULT current_timestamp(),
  `price_validity_period_end_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `purchase_order_lines`
--

LOCK TABLES `purchase_order_lines` WRITE;
/*!40000 ALTER TABLE `purchase_order_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `purchase_order_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `receipt_advice_headers`
--

DROP TABLE IF EXISTS `receipt_advice_headers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `receipt_advice_headers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `rcpth_id` varchar(100) DEFAULT '',
  `receipt_advice_type_code` varchar(50) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `line_count_numeric` int(10) unsigned DEFAULT 0,
  `order_id` int(10) unsigned DEFAULT 0,
  `despatch_id` int(10) unsigned DEFAULT 0,
  `delivery_customer_party_id` int(10) unsigned DEFAULT 0,
  `despatch_supplier_party_id` int(10) unsigned DEFAULT 0,
  `buyer_customer_party_id` int(10) unsigned DEFAULT 0,
  `seller_supplier_party_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `receipt_advice_headers`
--

LOCK TABLES `receipt_advice_headers` WRITE;
/*!40000 ALTER TABLE `receipt_advice_headers` DISABLE KEYS */;
/*!40000 ALTER TABLE `receipt_advice_headers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `receipt_advice_lines`
--

DROP TABLE IF EXISTS `receipt_advice_lines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `receipt_advice_lines` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `rcptl_id` varchar(100) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `received_quantity` int(10) unsigned DEFAULT 0,
  `short_quantity` int(10) unsigned DEFAULT 0,
  `shortage_action_code` varchar(50) DEFAULT '',
  `rejected_quantity` int(10) unsigned DEFAULT 0,
  `reject_reason_code` varchar(50) DEFAULT '',
  `reject_reason` varchar(50) DEFAULT '',
  `reject_action_code` varchar(50) DEFAULT '',
  `quantity_discrepancy_code` varchar(50) DEFAULT '',
  `oversupply_quantity` int(10) unsigned DEFAULT 0,
  `timing_complaint_code` varchar(50) DEFAULT '',
  `timing_complaint` varchar(50) DEFAULT '',
  `order_line_id` int(10) unsigned DEFAULT 0,
  `despatch_line_id` int(10) unsigned DEFAULT 0,
  `item_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `receipt_advice_header_id` int(10) unsigned DEFAULT 0,
  `received_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `receipt_advice_lines`
--

LOCK TABLES `receipt_advice_lines` WRITE;
/*!40000 ALTER TABLE `receipt_advice_lines` DISABLE KEYS */;
/*!40000 ALTER TABLE `receipt_advice_lines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `shipments`
--

DROP TABLE IF EXISTS `shipments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shipments` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `sh_id` varchar(100) DEFAULT '',
  `shipping_priority_level_code` varchar(100) DEFAULT '',
  `handling_code` varchar(100) DEFAULT '',
  `handling_instructions` varchar(50) DEFAULT '',
  `information` varchar(50) DEFAULT '',
  `gross_weight_measure` double DEFAULT 0,
  `net_weight_measure` double DEFAULT 0,
  `net_net_weight_measure` double DEFAULT 0,
  `gross_volume_measure` double DEFAULT 0,
  `net_volume_measure` double DEFAULT 0,
  `total_goods_item_quantity` bigint(20) DEFAULT 0,
  `total_transport_handling_unit_quantity` bigint(20) DEFAULT 0,
  `insurance_value_amount` double DEFAULT 0,
  `declared_customs_value_amount` double DEFAULT 0,
  `declared_for_carriage_value_amount` double DEFAULT 0,
  `declared_statistics_value_amount` double DEFAULT 0,
  `free_on_board_value_amount` double DEFAULT 0,
  `special_instructions` varchar(50) DEFAULT '',
  `delivery_instructions` varchar(50) DEFAULT '',
  `split_consignment_indicator` tinyint(1) DEFAULT 0,
  `consignment_quantity` bigint(20) DEFAULT 0,
  `return_address_id` int(10) unsigned DEFAULT 0,
  `origin_address_id` int(10) unsigned DEFAULT 0,
  `first_arrival_port_location_id` int(10) unsigned DEFAULT 0,
  `last_exit_port_location_id` int(10) unsigned DEFAULT 0,
  `export_country_id_code` varchar(50) DEFAULT '',
  `export_country_name` varchar(100) DEFAULT '',
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `shipments`
--

LOCK TABLES `shipments` WRITE;
/*!40000 ALTER TABLE `shipments` DISABLE KEYS */;
/*!40000 ALTER TABLE `shipments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tax_categories`
--

DROP TABLE IF EXISTS `tax_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tax_categories` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `tc_id` varchar(100) DEFAULT '',
  `tax_category_name` varchar(100) DEFAULT '',
  `percent` double DEFAULT 0,
  `base_unit_measure` varchar(20) DEFAULT '',
  `per_unit_amount` double DEFAULT 0,
  `tax_exemption_reason_code` varchar(50) DEFAULT '',
  `tax_exemption_reason` varchar(500) DEFAULT NULL,
  `tier_range` varchar(100) DEFAULT '',
  `tier_rate_percent` double DEFAULT 0,
  `tax_scheme_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tax_categories`
--

LOCK TABLES `tax_categories` WRITE;
/*!40000 ALTER TABLE `tax_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `tax_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tax_scheme_jurisdictions`
--

DROP TABLE IF EXISTS `tax_scheme_jurisdictions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tax_scheme_jurisdictions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `address_id` int(10) unsigned DEFAULT 0,
  `tax_scheme_id` int(10) unsigned DEFAULT 0,
  `tax_scheme_name` varchar(200) DEFAULT '',
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tax_scheme_jurisdictions`
--

LOCK TABLES `tax_scheme_jurisdictions` WRITE;
/*!40000 ALTER TABLE `tax_scheme_jurisdictions` DISABLE KEYS */;
/*!40000 ALTER TABLE `tax_scheme_jurisdictions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tax_schemes`
--

DROP TABLE IF EXISTS `tax_schemes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tax_schemes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `ts_id` varchar(100) DEFAULT '',
  `tax_scheme_name` varchar(100) DEFAULT '',
  `tax_type_code` varchar(20) DEFAULT '',
  `currency_code` varchar(20) DEFAULT '',
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tax_schemes`
--

LOCK TABLES `tax_schemes` WRITE;
/*!40000 ALTER TABLE `tax_schemes` DISABLE KEYS */;
/*!40000 ALTER TABLE `tax_schemes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tax_sub_totals`
--

DROP TABLE IF EXISTS `tax_sub_totals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tax_sub_totals` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `taxable_amount` double DEFAULT 0,
  `tax_amount` double DEFAULT 0,
  `calculation_sequence_numeric` int(10) unsigned DEFAULT 0,
  `transaction_currency_tax_amount` double DEFAULT 0,
  `percent` double DEFAULT 0,
  `base_unit_measure` varchar(20) DEFAULT '',
  `per_unit_amount` double DEFAULT 0,
  `tier_range` varchar(50) DEFAULT '',
  `tier_rate_percent` double DEFAULT 0,
  `tax_category_id` int(10) unsigned DEFAULT 0,
  `tax_total_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tax_sub_totals`
--

LOCK TABLES `tax_sub_totals` WRITE;
/*!40000 ALTER TABLE `tax_sub_totals` DISABLE KEYS */;
/*!40000 ALTER TABLE `tax_sub_totals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tax_totals`
--

DROP TABLE IF EXISTS `tax_totals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tax_totals` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `tax_amount` double DEFAULT 0,
  `rounding_amount` double DEFAULT 0,
  `tax_evidence_indicator` tinyint(1) DEFAULT 0,
  `tax_included_indicator` tinyint(1) DEFAULT 0,
  `master_flag` varchar(20) DEFAULT '',
  `master_id` int(10) unsigned DEFAULT 0,
  `tax_category_id` int(10) unsigned DEFAULT 0,
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tax_totals`
--

LOCK TABLES `tax_totals` WRITE;
/*!40000 ALTER TABLE `tax_totals` DISABLE KEYS */;
/*!40000 ALTER TABLE `tax_totals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `waybills`
--

DROP TABLE IF EXISTS `waybills`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `waybills` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid4` binary(16) DEFAULT NULL,
  `waybill_id` varchar(100) DEFAULT '',
  `carrier_assigned_id` varchar(100) DEFAULT '',
  `name1` varchar(50) DEFAULT '',
  `description` varchar(50) DEFAULT '',
  `note` varchar(50) DEFAULT '',
  `shipping_order_id` varchar(100) DEFAULT '',
  `ad_valorem_indicator` tinyint(1) DEFAULT 0,
  `declared_carriage_value_amount` double DEFAULT 0,
  `declared_carriage_value_amount_currency_code` varchar(20) DEFAULT '',
  `other_instruction` varchar(50) DEFAULT '',
  `consignor_party_id` int(10) unsigned DEFAULT 0,
  `carrier_party_id` int(10) unsigned DEFAULT 0,
  `freight_forwarder_party_id` int(10) unsigned DEFAULT 0,
  `shipment_id` int(10) unsigned DEFAULT 0,
  `issue_date` datetime DEFAULT current_timestamp(),
  `status_code` varchar(50) DEFAULT 'active',
  `created_by_user_id` varchar(50) DEFAULT 'active',
  `updated_by_user_id` varchar(50) DEFAULT 'active',
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `waybills`
--

LOCK TABLES `waybills` WRITE;
/*!40000 ALTER TABLE `waybills` DISABLE KEYS */;
/*!40000 ALTER TABLE `waybills` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-10-22 15:03:24
