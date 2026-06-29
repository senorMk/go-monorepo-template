import 'dart:async';

import 'package:flutter_bloc/flutter_bloc.dart';

import '../domain/entities/app_session.dart';
import '../domain/entities/app_user.dart';
import '../domain/repositories/auth_repository.dart';

class AuthState {
  const AuthState({
    this.user,
    this.isLoading = false,
    this.error,
    this.message,
  });

  final AppUser? user;
  final bool isLoading;
  final String? error;
  final String? message;

  bool get isAuthenticated => user != null;

  AuthState copyWith({
    AppUser? user,
    bool? isLoading,
    String? error,
    String? message,
    bool clearError = false,
    bool clearMessage = false,
  }) {
    return AuthState(
      user: user ?? this.user,
      isLoading: isLoading ?? this.isLoading,
      error: clearError ? null : (error ?? this.error),
      message: clearMessage ? null : (message ?? this.message),
    );
  }
}

class AuthCubit extends Cubit<AuthState> {
  AuthCubit(this._authRepository) : super(const AuthState()) {
    _subscription = _authRepository.watchAuthState().listen((AppUser? user) {
      emit(state.copyWith(user: user, isLoading: false, clearError: true));
    });
  }

  final AuthRepository _authRepository;
  StreamSubscription<AppUser?>? _subscription;

  Future<void> bootstrap() async {
    emit(state.copyWith(isLoading: true, clearError: true, clearMessage: true));
    final user = await _authRepository.getCurrentUser();
    emit(state.copyWith(user: user, isLoading: false));
  }

  Future<void> signIn(String email, String password) async {
    emit(state.copyWith(isLoading: true, clearError: true));
    try {
      final session = await _authRepository.signInWithEmail(email: email, password: password);
      emit(state.copyWith(user: session.user, isLoading: false));
    } catch (e) {
      emit(state.copyWith(isLoading: false, error: e.toString()));
    }
  }

  Future<void> signOut() async {
    await _authRepository.signOut();
    emit(const AuthState());
  }

  @override
  Future<void> close() {
    _subscription?.cancel();
    return super.close();
  }
}
