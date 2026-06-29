import '../../domain/entities/app_session.dart';
import '../../domain/entities/app_user.dart';
import '../../domain/repositories/auth_repository.dart';
import 'auth_exceptions.dart';

class UnimplementedAuthRepository implements AuthRepository {
  const UnimplementedAuthRepository();

  @override
  Future<AppUser?> getCurrentUser() async {
    return null;
  }

  @override
  Future<AppSession> signInWithEmail({
    required String email,
    required String password,
  }) {
    throw const AuthNotConfiguredException();
  }

  @override
  Future<AppSession> signUpWithEmail({
    required String email,
    required String password,
    String? displayName,
  }) {
    throw const AuthNotConfiguredException();
  }

  @override
  Future<void> signOut() {
    throw const AuthNotConfiguredException();
  }

  @override
  Future<void> sendPasswordResetEmail({required String email}) {
    throw const AuthNotConfiguredException();
  }

  @override
  Stream<AppUser?> watchAuthState() {
    return const Stream<AppUser?>.empty();
  }
}
