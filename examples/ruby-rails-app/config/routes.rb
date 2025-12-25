Rails.application.routes.draw do
  namespace :api do
    namespace :v1 do
      resources :users do
        resources :posts, only: [:index, :create]
      end
      resources :posts, only: [:index, :show, :update, :destroy]
    end
  end

  root 'home#index'
  get 'health', to: 'health#check'
end
