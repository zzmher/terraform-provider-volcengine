resource "volcengine_tls_index" "foo" {
  topic_id = "65d67d34-c5b4-4ec8-b3a9-175d3366****"

  full_text {
    case_sensitive = true
    delimiter = "!"
    include_chinese = false
  }

  key_value {
    key = "k1"
    value_type = "json"
    case_sensitive = true
    delimiter = "!"
    include_chinese = false
    sql_flag = false
    json_keys {
      key = "k2.k4"
      value_type = "text"
    }
    json_keys {
      key = "k3.k4"
      value_type = "long"
    }
  }

  key_value {
    key = "k5"
    value_type = "text"
    case_sensitive = true
    delimiter = "!"
    include_chinese = false
    sql_flag = false
  }
}