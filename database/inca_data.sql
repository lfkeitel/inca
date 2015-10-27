-- phpMyAdmin SQL Dump
-- version 4.0.10deb1
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: Oct 26, 2015 at 09:24 PM
-- Server version: 5.5.44-0ubuntu0.14.04.1
-- PHP Version: 5.5.9-1ubuntu4.13

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Database: `inca_data`
--

-- --------------------------------------------------------

--
-- Table structure for table `configs`
--

CREATE TABLE IF NOT EXISTS `configs` (
  `configid` int(11) NOT NULL AUTO_INCREMENT,
  `deviceid` int(11) NOT NULL,
  `timestamp` int(11) NOT NULL,
  `filename` text NOT NULL,
  `is_ok` tinyint(1) NOT NULL,
  `parsed_config` longtext NOT NULL,
  PRIMARY KEY (`configid`),
  KEY `deviceid` (`deviceid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=7 ;

-- --------------------------------------------------------

--
-- Table structure for table `conn_profiles`
--

CREATE TABLE IF NOT EXISTS `conn_profiles` (
  `profileid` int(11) NOT NULL AUTO_INCREMENT,
  `name` tinytext NOT NULL,
  `protocol` tinytext NOT NULL,
  `username` text NOT NULL,
  `password` text NOT NULL,
  `enable` text NOT NULL,
  PRIMARY KEY (`profileid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=5 ;

-- --------------------------------------------------------

--
-- Table structure for table `devices`
--

CREATE TABLE IF NOT EXISTS `devices` (
  `deviceid` int(11) NOT NULL AUTO_INCREMENT,
  `name` tinytext NOT NULL,
  `hostname` tinytext NOT NULL,
  `conn_profile` int(11) NOT NULL,
  `manufacturer` tinytext NOT NULL,
  `model` tinytext NOT NULL,
  `custom` tinyint(1) NOT NULL DEFAULT '1',
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `parse_config` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`deviceid`),
  KEY `conn_profile` (`conn_profile`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=16 ;

-- --------------------------------------------------------

--
-- Table structure for table `device_status`
--

CREATE TABLE IF NOT EXISTS `device_status` (
  `statusid` int(11) NOT NULL AUTO_INCREMENT,
  `deviceid` int(11) NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `last_polled` int(11) NOT NULL,
  `last_error` text NOT NULL,
  PRIMARY KEY (`statusid`),
  KEY `deviceid` (`deviceid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=14 ;

-- --------------------------------------------------------

--
-- Table structure for table `jobs`
--

CREATE TABLE IF NOT EXISTS `jobs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `status` tinytext NOT NULL,
  `last_error` text NOT NULL,
  `start_time` int(11) NOT NULL,
  `end_time` int(11) NOT NULL,
  `run_time` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

-- --------------------------------------------------------

--
-- Table structure for table `logs`
--

CREATE TABLE IF NOT EXISTS `logs` (
  `logid` int(11) NOT NULL AUTO_INCREMENT,
  `area` tinytext NOT NULL,
  `level` tinyint(4) NOT NULL,
  `timestamp` bigint(20) NOT NULL,
  `message` text NOT NULL,
  PRIMARY KEY (`logid`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=78 ;

-- --------------------------------------------------------

--
-- Table structure for table `scripts`
--

CREATE TABLE IF NOT EXISTS `scripts` (
  `scriptid` int(11) NOT NULL AUTO_INCREMENT,
  `name` tinytext NOT NULL,
  `last_edited` int(11) NOT NULL,
  `script_text` longtext NOT NULL,
  PRIMARY KEY (`scriptid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `configs`
--
ALTER TABLE `configs`
  ADD CONSTRAINT `device_id_fk1` FOREIGN KEY (`deviceid`) REFERENCES `devices` (`deviceid`) ON DELETE CASCADE;

--
-- Constraints for table `devices`
--
ALTER TABLE `devices`
  ADD CONSTRAINT `conn_profile_fk1` FOREIGN KEY (`conn_profile`) REFERENCES `conn_profiles` (`profileid`);

--
-- Constraints for table `device_status`
--
ALTER TABLE `device_status`
  ADD CONSTRAINT `device_status_ibfk_1` FOREIGN KEY (`deviceid`) REFERENCES `devices` (`deviceid`) ON DELETE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
