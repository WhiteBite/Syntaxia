import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/services/api_service.dart';
import '../models/home_data.dart';

final homeProvider = FutureProvider<HomeData>((ref) async {
  final apiService = ref.watch(apiServiceProvider);
  
  // Simulate API call
  await Future.delayed(const Duration(seconds: 1));
  
  return HomeData(
    userCount: 1234,
    activeUsers: 567,
    totalRevenue: 89012.34,
    lastUpdated: DateTime.now(),
  );
});

final apiServiceProvider = Provider<ApiService>((ref) {
  return ApiService();
});