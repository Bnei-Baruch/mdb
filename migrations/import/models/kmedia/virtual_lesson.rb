class VirtualLesson < ActiveRecord::Base
  establish_connection $kmedia_config

  has_many :containers
end
