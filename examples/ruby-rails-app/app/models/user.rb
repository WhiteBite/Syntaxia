class User < ApplicationRecord
  has_many :posts, dependent: :destroy
  has_secure_password

  validates :name, presence: true, length: { minimum: 2, maximum: 100 }
  validates :email, presence: true, uniqueness: { case_sensitive: false },
                    format: { with: URI::MailTo::EMAIL_REGEXP }

  before_save :downcase_email

  scope :active, -> { where(active: true) }
  scope :recent, -> { order(created_at: :desc) }

  def full_info
    "#{name} (#{email})"
  end

  private

  def downcase_email
    self.email = email.downcase
  end
end
