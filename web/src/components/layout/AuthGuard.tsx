import { useAuthStore } from "@/stores";
import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";

export default function AuthGuard() {
	const navigate = useNavigate();
	const { isAuthenticated } = useAuthStore();

	useEffect(() => {
		if (!isAuthenticated) {
			navigate("/auth/login", { replace: true });
		}
	}, [isAuthenticated, navigate]);

	if (!isAuthenticated) {
		return null;
	}

	return <Outlet />;
}
