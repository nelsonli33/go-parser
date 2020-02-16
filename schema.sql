CREATE TABLE traffic_accident (
	id INT PRIMARY KEY AUTO_INCREMENT,
  `date` DATETIME,
  death_count INT,
  injury_count INT,
  latitude FLOAT4,
  longitude FLOAT4
)