class UserService
  def initialize(user = nil)
    @user = user
  end

  def create(params)
    user = User.new(params)
    if user.save
      send_welcome_email(user)
      { success: true, user: user }
    else
      { success: false, errors: user.errors.full_messages }
    end
  end

  def update(params)
    if @user.update(params)
      { success: true, user: @user }
    else
      { success: false, errors: @user.errors.full_messages }
    end
  end

  def deactivate
    @user.update(active: false)
    send_deactivation_email(@user)
    { success: true, message: 'User deactivated' }
  end

  def generate_auth_token
    token = SecureRandom.hex(32)
    @user.update(auth_token: token)
    token
  end

  private

  def send_welcome_email(user)
    # UserMailer.welcome(user).deliver_later
    Rails.logger.info("Welcome email sent to #{user.email}")
  end

  def send_deactivation_email(user)
    # UserMailer.deactivation(user).deliver_later
    Rails.logger.info("Deactivation email sent to #{user.email}")
  end
end
