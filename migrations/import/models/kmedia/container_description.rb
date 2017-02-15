class ContainerDescription < ActiveRecord::Base
  establish_connection $kmedia_config
  belongs_to :container
end
