import 'dart:async';

import 'package:APP_SNAKE/cubit/auth_cubit.dart';
import 'package:APP_SNAKE/domain/entities/app_session.dart';
import 'package:APP_SNAKE/domain/entities/app_user.dart';
import 'package:APP_SNAKE/domain/repositories/auth_repository.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('sign in exposes the session user and sign out clears it', () async {
    final repository = _FakeAuthRepository();
    final cubit = AuthCubit(repository);
    addTearDown(() async {
      await cubit.close();
      await repository.close();
    });

    await cubit.bootstrap();
    expect(cubit.state.isAuthenticated, isFalse);

    await cubit.signIn('user@example.com', 'password');
    expect(cubit.state.user, repository.user);
    expect(cubit.state.isLoading, isFalse);

    await cubit.signOut();
    expect(cubit.state.isAuthenticated, isFalse);
  });
}

class _FakeAuthRepository implements AuthRepository {
  final user = const AppUser(id: '1', email: 'user@example.com');
  final _changes = StreamController<AppUser?>.broadcast();

  Future<void> close() => _changes.close();

  @override
  Future<AppUser?> getCurrentUser() async => null;

  @override
  Stream<AppUser?> watchAuthState() => _changes.stream;

  @override
  Future<AppSession> signInWithEmail({required String email, required String password}) async {
    return AppSession(user: user);
  }

  @override
  Future<void> signOut() async {}

  @override
  Future<AppSession> signUpWithEmail({required String email, required String password, String? displayName}) {
    throw UnimplementedError();
  }

  @override
  Future<void> sendPasswordResetEmail({required String email}) {
    throw UnimplementedError();
  }
}
