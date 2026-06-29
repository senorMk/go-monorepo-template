import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:APP_SNAKE/cubit/auth_cubit.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AuthCubit, AuthState>(
      builder: (context, state) {
        return Scaffold(
          appBar: AppBar(title: const Text('APP_DISPLAY_NAME')),
          body: Center(
            child: state.isAuthenticated
                ? Text('Welcome, ${state.user!.email}')
                : const Text('Not signed in'),
          ),
        );
      },
    );
  }
}
