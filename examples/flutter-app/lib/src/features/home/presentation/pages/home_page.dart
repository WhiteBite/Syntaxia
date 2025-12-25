import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../shared/widgets/app_bar.dart';
import '../../../shared/widgets/bottom_nav_bar.dart';
import '../providers/home_provider.dart';

class HomePage extends ConsumerWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final homeState = ref.watch(homeProvider);
    
    return Scaffold(
      appBar: const CustomAppBar(title: 'Home'),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(homeProvider);
        },
        child: homeState.when(
          data: (data) => _buildContent(context, data),
          loading: () => const Center(
            child: CircularProgressIndicator(),
          ),
          error: (error, stack) => Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(
                  Icons.error_outline,
                  size: 64,
                  color: Colors.red,
                ),
                const SizedBox(height: 16),
                Text(
                  'Error: $error',
                  style: Theme.of(context).textTheme.bodyLarge,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () => ref.invalidate(homeProvider),
                  child: const Text('Retry'),
                ),
              ],
            ),
          ),
        ),
      ),
      bottomNavigationBar: const CustomBottomNavBar(currentIndex: 0),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          context.go('/profile');
        },
        child: const Icon(Icons.person),
      ),
    );
  }

  Widget _buildContent(BuildContext context, HomeData data) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Welcome!',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 8),
                Text(
                  'This is a test Flutter app for SYNTAXIA.',
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: ListTile(
            leading: const Icon(Icons.analytics),
            title: const Text('Analytics'),
            subtitle: Text('Users: ${data.userCount}'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              // Navigate to analytics
            },
          ),
        ),
        Card(
          child: ListTile(
            leading: const Icon(Icons.settings),
            title: const Text('Settings'),
            subtitle: const Text('App configuration'),
            trailing: const Icon(Icons.arrow_forward_ios),
            onTap: () {
              context.go('/settings');
            },
          ),
        ),
      ],
    );
  }
}