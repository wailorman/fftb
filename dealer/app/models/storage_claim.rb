class StorageClaim < ApplicationRecord
  DEFAULT_URL_TTL = 1.day

  enum kind: { s3: 's3' }, _suffix: true

  enum purpose: { none: 'none',
                  convert_input: 'convert_input',
                  convert_output: 'convert_output' }, _suffix: true

  belongs_to :task

  validates :path, presence: true
  validates :purpose, presence: true
  validates :type, inclusion: { in: %w[InputStorageClaim OutputStorageClaim] }
end
