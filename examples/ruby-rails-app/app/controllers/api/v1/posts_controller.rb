module Api
  module V1
    class PostsController < ApplicationController
      before_action :set_post, only: [:show, :update, :destroy]
      before_action :set_user, only: [:index, :create], if: -> { params[:user_id].present? }

      def index
        @posts = if @user
                   @user.posts.recent
                 else
                   Post.published.recent.page(params[:page]).per(20)
                 end
        render json: @posts
      end

      def show
        render json: @post, include: :user
      end

      def create
        @post = @user.posts.build(post_params)
        if @post.save
          render json: @post, status: :created
        else
          render json: { errors: @post.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def update
        if @post.update(post_params)
          render json: @post
        else
          render json: { errors: @post.errors.full_messages }, status: :unprocessable_entity
        end
      end

      def destroy
        @post.destroy
        head :no_content
      end

      private

      def set_post
        @post = Post.find(params[:id])
      end

      def set_user
        @user = User.find(params[:user_id])
      end

      def post_params
        params.require(:post).permit(:title, :content, :published)
      end
    end
  end
end
