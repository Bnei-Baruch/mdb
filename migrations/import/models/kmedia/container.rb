class Container < ActiveRecord::Base
  establish_connection $kmedia_config
end
