import '../entities/app_session.dart';
import '../entities/app_user.dart';

abstract class AuthRepository {
  Future<AppSession> signInWithEmail({
    required String email,
    required String password,
  });

  Future<AppSession> signUpWithEmail({
    required String email,
    required String password,
    String? displayName,
  });

  Future<void> signOut();

  Future<void> sendPasswordResetEmail({required String email});

  Future<AppUser?> getCurrentUser();

  Stream<AppUser?> watchAuthState();
}
