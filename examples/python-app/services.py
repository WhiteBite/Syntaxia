from models import User, Post
from typing import List, Optional, Dict, Any

class UserService:
    def __init__(self, db):
        self.db = db
    
    def get_all_users(self) -> List[User]:
        return User.query.all()
    
    def get_user_by_id(self, user_id: int) -> Optional[User]:
        return User.query.get(user_id)
    
    def create_user(self, data: Dict[str, Any]) -> User:
        user = User(name=data['name'], email=data['email'])
        self.db.session.add(user)
        self.db.session.commit()
        return user
    
    def update_user(self, user_id: int, data: Dict[str, Any]) -> Optional[User]:
        user = User.query.get(user_id)
        if user:
            user.name = data.get('name', user.name)
            user.email = data.get('email', user.email)
            self.db.session.commit()
        return user
    
    def delete_user(self, user_id: int) -> bool:
        user = User.query.get(user_id)
        if user:
            self.db.session.delete(user)
            self.db.session.commit()
            return True
        return False

class PostService:
    def __init__(self, db):
        self.db = db
    
    def get_all_posts(self) -> List[Post]:
        return Post.query.all()
    
    def get_posts_by_user(self, user_id: int) -> List[Post]:
        return Post.query.filter_by(user_id=user_id).all()
    
    def create_post(self, data: Dict[str, Any]) -> Post:
        post = Post(
            title=data['title'],
            content=data['content'],
            user_id=data['user_id']
        )
        self.db.session.add(post)
        self.db.session.commit()
        return post
