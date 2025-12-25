class HomeData {
  final int userCount;
  final int activeUsers;
  final double totalRevenue;
  final DateTime lastUpdated;

  const HomeData({
    required this.userCount,
    required this.activeUsers,
    required this.totalRevenue,
    required this.lastUpdated,
  });

  HomeData copyWith({
    int? userCount,
    int? activeUsers,
    double? totalRevenue,
    DateTime? lastUpdated,
  }) {
    return HomeData(
      userCount: userCount ?? this.userCount,
      activeUsers: activeUsers ?? this.activeUsers,
      totalRevenue: totalRevenue ?? this.totalRevenue,
      lastUpdated: lastUpdated ?? this.lastUpdated,
    );
  }

  @override
  String toString() {
    return 'HomeData(userCount: $userCount, activeUsers: $activeUsers, totalRevenue: $totalRevenue, lastUpdated: $lastUpdated)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    
    return other is HomeData &&
        other.userCount == userCount &&
        other.activeUsers == activeUsers &&
        other.totalRevenue == totalRevenue &&
        other.lastUpdated == lastUpdated;
  }

  @override
  int get hashCode {
    return userCount.hashCode ^
        activeUsers.hashCode ^
        totalRevenue.hashCode ^
        lastUpdated.hashCode;
  }
}