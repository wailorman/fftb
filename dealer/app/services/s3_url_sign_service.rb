class S3UrlSignService < ApplicationService
  attr_reader :path,
              :provider

  def initialize(storage_claim)
    @provider = storage_claim.provider
    @path = storage_claim.path
    super()
  end

  def put(expires_in: nil)
    object.presigned_url(:put, expires_in: expires_in&.to_i)
  end

  def get(expires_in: nil)
    object.presigned_url(:get, expires_in: expires_in&.to_i)
  end

  private def s3_config
    @s3_config ||= Rails.application.config.application_options[:s3][provider.to_sym]
  end

  private def aws_client
    @aws_client ||= Aws::S3::Client.new(s3_config.except(:bucket))
  end

  private def s3
    @s3 ||= Aws::S3::Resource.new(client: aws_client)
  end

  private def bucket
    @bucket ||= s3.bucket(s3_config[:bucket])
  end

  private def object
    @object ||= bucket.object(path)
  end
end
