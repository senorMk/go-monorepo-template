import 'dart:io';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_displaymode/flutter_displaymode.dart';
import 'package:APP_SNAKE/config/style.dart';
import 'package:APP_SNAKE/cubit/auth_cubit.dart';
import 'package:APP_SNAKE/cubit/theme_cubit.dart';
import 'package:APP_SNAKE/data/auth/unimplemented_auth_repository.dart';
import 'package:APP_SNAKE/domain/repositories/auth_repository.dart';
import 'package:APP_SNAKE/ui/screens/home_screen.dart';
import 'package:hydrated_bloc/hydrated_bloc.dart';
import 'package:path_provider/path_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await EasyLocalization.ensureInitialized();
  if (!kIsWeb && Platform.isAndroid) {
    await FlutterDisplayMode.setHighRefreshRate();
  }

  final appDir = await getApplicationDocumentsDirectory();
  final tmpDir = await getTemporaryDirectory();
  HydratedBloc.storage = await HydratedStorage.build(storageDirectory: tmpDir);

  runApp(
    EasyLocalization(
      path: 'assets/translations',
      supportedLocales: const [Locale('en')],
      fallbackLocale: const Locale('en'),
      useFallbackTranslations: true,
      child: const MyApp(authRepository: UnimplementedAuthRepository()),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key, required this.authRepository});

  final AuthRepository authRepository;

  @override
  Widget build(BuildContext context) {
    return BlocProvider<AuthCubit>(
      create: (_) => AuthCubit(authRepository)..bootstrap(),
      child: BlocProvider<ThemeCubit>(
        create: (_) => ThemeCubit(),
        child: BlocBuilder<ThemeCubit, ThemeModeState>(
          builder: (context, state) {
            return MaterialApp(
              title: 'APP_DISPLAY_NAME',
              theme: Style.lightTheme,
              darkTheme: Style.darkTheme,
              themeMode: state.themeMode,
              localizationsDelegates: context.localizationDelegates,
              supportedLocales: context.supportedLocales,
              locale: context.locale,
              debugShowCheckedModeBanner: false,
              home: const HomeScreen(),
            );
          },
        ),
      ),
    );
  }
}
