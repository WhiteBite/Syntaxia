module Api
  module V1
    class UsersController < ApplicationController
      skip_before_action :authenticate_request, only: [:create]

      before_action :set_user, only: [:show, :update, :destroy]

      def index
        @users = User.active.recent.page(params[:page]).per(20)
        render json: @users
      end

      def show
        render json: @user, include: :posts
      end

      def create
        @user = User.new(user_params)
        if @user.save
          render json: @user, status: :created
        else
          render json: { errors: @user.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def update
        if @user.update(user_params)
          render json: @user
        else
          render json: { errors: @user.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        @user.destroy
        head :no_content
      end

      private

      def set_user
        @user = User.find(params[:id])
      end

      def user_params
        params.require(:user).permit(:name, :email, :password, :password_confirmation)
      end
    end
  end
end
