class AuthNotConfiguredException implements Exception {
  const AuthNotConfiguredException();

  @override
  String toString() =>
      'Auth backend not configured yet. Add a concrete provider implementation.';
}
