import 'package:equatable/equatable.dart';

import 'app_user.dart';

class AppSession extends Equatable {
  const AppSession({
    this.user,
    this.accessToken,
    this.refreshToken,
    this.expiresAt,
    this.requiresEmailConfirmation = false,
  });

  final AppUser? user;
  final String? accessToken;
  final String? refreshToken;
  final DateTime? expiresAt;
  final bool requiresEmailConfirmation;

  @override
  List<Object?> get props => <Object?>[
        user,
        accessToken,
        refreshToken,
        expiresAt,
        requiresEmailConfirmation,
      ];
}
