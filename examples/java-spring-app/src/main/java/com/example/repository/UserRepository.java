package com.example.repository;

import com.example.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;

/**
 * Repository interface for User entity
 */
@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    
    Optional<User> findByEmail(String email);
    
    boolean existsByEmail(String email);
    
    @Query("SELECT u FROM User u WHERE u.name LIKE %?1%")
    java.util.List<User> findByNameContaining(String name);
    
    @Query("SELECT COUNT(u) FROM User u WHERE u.createdAt >= ?1")
    long countUsersCreatedAfter(java.time.LocalDateTime date);
}