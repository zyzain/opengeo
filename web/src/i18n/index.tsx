import { type ReactNode, createContext, useContext, useState } from "react";
import { IntlProvider } from "react-intl";
import enUS from "./en-US.json";
import zhCN from "./zh-CN.json";

export type Locale = "zh-CN" | "en-US";

const messages: Record<Locale, Record<string, string>> = {
	"zh-CN": zhCN,
	"en-US": enUS,
};

interface I18nContextType {
	locale: Locale;
	setLocale: (locale: Locale) => void;
}

const I18nContext = createContext<I18nContextType>({
	locale: "zh-CN",
	setLocale: () => {},
});

export function useI18n() {
	return useContext(I18nContext);
}

interface I18nProviderProps {
	children: ReactNode;
}

export function I18nProvider({ children }: I18nProviderProps) {
	const [locale, setLocale] = useState<Locale>(() => {
		if (typeof window !== "undefined") {
			return (localStorage.getItem("locale") as Locale) || "zh-CN";
		}
		return "zh-CN";
	});

	const handleSetLocale = (newLocale: Locale) => {
		setLocale(newLocale);
		localStorage.setItem("locale", newLocale);
	};

	return (
		<I18nContext.Provider value={{ locale, setLocale: handleSetLocale }}>
			<IntlProvider locale={locale} messages={messages[locale]}>
				{children}
			</IntlProvider>
		</I18nContext.Provider>
	);
}
