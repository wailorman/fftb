class StorageClaim < ApplicationRecord
  DEFAULT_URL_TTL = 1.day

  enum kind: { s3: 's3' }

  validates :path, presence: true
end
