require 'active_record'
require 'yaml'
$kmedia_config = YAML::load(File.open('config/database.yml'))['kmedia']
$mdb_config = YAML::load(File.open('config/database.yml'))['mdb']
Dir[File.dirname(__FILE__) + '/models/**/*.rb'].each {|file| require file }

# virual lesson (virtual_lessons) -> collection
# lesson part (containers) -> content unit
# files (file_assets) -> files


puts VirtualLesson.count
puts FileAsset.count
puts Container.count

puts Collection.count
puts ContentUnit.count
puts MDBFile.count