FactoryBot.define do
  factory :storage_claim, class: 'InputStorageClaim' do
    type { 'InputStorageClaim' }
    kind { :s3 }
    provider { :yandex }
    path { "claims/#{SecureRandom.uuid}/input.mp4" }
    name { 'input.mp4' }
    purpose { :convert_input }

    factory :input_storage_claim, class: 'InputStorageClaim' do
    end

    factory :output_storage_claim, class: 'OutputStorageClaim' do
      name { 'output.mp4' }
      purpose { :convert_output }
    end
  end
end
