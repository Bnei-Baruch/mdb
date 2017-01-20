require 'active_record'
require 'yaml'


# virual lesson (virtual_lessons) -> collection
# lesson part (containers) -> content unit
# files (file_assets) -> files

@db_config = YAML::load(File.open('config/database.yml'))

class VirtualLesson < ActiveRecord::Base
  establish_connection @db_config['kmedia']
end

puts VirtualLesson.count


