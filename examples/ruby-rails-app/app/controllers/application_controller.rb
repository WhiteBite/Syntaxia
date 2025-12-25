class ApplicationController < ActionController::API
  include ActionController::HttpAuthentication::Token::ControllerMethods

  before_action :authenticate_request

  rescue_from ActiveRecord::RecordNotFound, with: :not_found
  rescue_from ActiveRecord::RecordInvalid, with: :unprocessable_entity

  private

  def authenticate_request
    authenticate_or_request_with_http_token do |token, _options|
      @current_user = User.find_by(auth_token: token)
    end
  end

  def current_user
    @current_user
  end

  def not_found
    render json: { error: 'Not found' }, status: :not_found
  end

  def unprocessable_entity(exception)
    render json: { errors: exception.record.errors.full_messages }, status: :unprocessable_entity
  end
end
