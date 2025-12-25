using csharp_app.Models;

namespace csharp_app.Services
{
    public interface IPostService
    {
        Task<IEnumerable<Post>> GetAllPostsAsync();
        Task<IEnumerable<Post>> GetPostsByUserIdAsync(int userId);
        Task<Post?> GetPostByIdAsync(int id);
        Task<Post> CreatePostAsync(Post post);
        Task<bool> DeletePostAsync(int id);
    }
}
