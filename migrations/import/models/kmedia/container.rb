class Container < ActiveRecord::Base
  establish_connection $kmedia_config
  belongs_to :virtual_lesson
  has_many :container_descriptions

  has_and_belongs_to_many :file_assets

  def name(lang)
    res = self.container_descriptions.where(lang: lang).first

    res.try(:container_desc)
  end
  def description(lang)
    res = self.container_descriptions.where(lang: lang).first

    res.try(:descr)
  end
end
