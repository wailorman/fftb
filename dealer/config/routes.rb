Rails.application.routes.draw do
  use_twirp

  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  root 'home#index'

  resources :home, only: %i[index]

  resources :remotes, only: %i[index], param: :name do
    resources :files, only: %i[index], controller: 'remotes/files'
  end
  resources :file_selections, only: %i[create show destroy] do
    resources :items, only: %i[], controller: 'file_selections/items' do
      member do
        post :hide
        post :reveal
      end
    end
  end
  resources :media_meta_reports, only: %i[show create]
  resources :orders, only: %i[index new create edit update show] do
    collection do
      get :choose_type
    end

    member do
      post :cancel
    end
  end
end
