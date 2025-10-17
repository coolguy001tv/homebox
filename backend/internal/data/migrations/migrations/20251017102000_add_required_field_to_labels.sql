-- add "required" field to "labels" table
ALTER TABLE `labels` ADD COLUMN `required` bool NOT NULL DEFAULT false;
