class FileAsset < ActiveRecord::Base
  establish_connection $kmedia_config

  has_and_belongs_to_many :containers
  has_many :file_asset_descriptions, foreign_key: :file_id


  def description(lang)
    self.file_asset_descriptions.where(lang: lang).try(:first).try(:filedesc)
  end
end