ALTER TABLE _ps_catalog.promscale_instance_information
  RENAME COLUMN promscale_metrics_queries_executed_total
  TO promscale_metrics_queries_success_total;