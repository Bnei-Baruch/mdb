class FileAssetDescription < ActiveRecord::Base
  establish_connection $kmedia_config

  belongs_to :file_asset, foreign_key: :file_id
end