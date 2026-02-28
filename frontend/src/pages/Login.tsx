import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { LogIn, Mail, Lock, AlertCircle, Loader2, Grid3X3 } from 'lucide-react';
import api from '../api/client';
import type { ApiResponse } from '../types';

interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', {
        email,
        password,
      });

      if (response.data.success) {
        localStorage.setItem('access_token', response.data.data.access_token);
        localStorage.setItem('refresh_token', response.data.data.refresh_token);
        navigate('/');
      } else {
        setError(response.data.message || '登录失败');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || '无法连接到服务器，请检查网络设置。');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#242424] text-white flex flex-col items-center justify-center p-4">
      <Link to="/" className="mb-8 flex items-center gap-3 text-3xl font-bold text-white hover:opacity-80 transition-opacity">
        <Grid3X3 className="w-10 h-10 text-[#55ff55]" />
        <span>YSM 模型站</span>
      </Link>

      <div className="modal_area w-full max-w-md bg-[#3c3c3c] border border-[#242424] p-8 shadow-2xl relative">
        <div className="flex flex-col items-center mb-8">
          <div className="w-16 h-16 rounded-full bg-[#242424] border border-[#3c3c3c] flex items-center justify-center mb-4">
            <LogIn className="w-8 h-8 text-[#55ff55]" />
          </div>
          <h2 className="text-2xl font-bold tracking-tight">欢迎回来</h2>
          <p className="text-sm text-[#aaaaaa] mt-2">请登录您的账户以管理模型</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/50 rounded flex items-start gap-3 text-red-500 text-sm">
            <AlertCircle className="w-5 h-5 shrink-0" />
            <p>{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label className="text-xs font-bold uppercase tracking-wider text-[#aaaaaa] flex items-center gap-2">
              <Mail className="w-3.5 h-3.5" />
              电子邮箱
            </label>
            <input
              type="email"
              required
              placeholder="your@email.com"
              className="input w-full bg-[#242424] border-[#3c3c3c] h-12 px-4 focus:border-[#55ff55] transition-all"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-xs font-bold uppercase tracking-wider text-[#aaaaaa] flex items-center gap-2">
                <Lock className="w-3.5 h-3.5" />
                密码
              </label>
              <Link to="/forgot-password" className="text-[10px] font-bold text-[#55ff55] hover:underline uppercase">
                忘记密码?
              </Link>
            </div>
            <input
              type="password"
              required
              placeholder="••••••••"
              className="input w-full bg-[#242424] border-[#3c3c3c] h-12 px-4 focus:border-[#55ff55] transition-all"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="green_btn w-full h-12 flex items-center justify-center gap-2 font-bold uppercase tracking-widest text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <>
                <Loader2 className="w-5 h-5 animate-spin" />
                登录中...
              </>
            ) : (
              '立即登录'
            )}
          </button>
        </form>

        <div className="mt-8 pt-6 border-t border-[#242424] text-center">
          <p className="text-sm text-[#aaaaaa]">
            还没有账号?{' '}
            <Link to="/register" className="text-[#55ff55] font-bold hover:underline">
              立即注册
            </Link>
          </p>
        </div>
      </div>

      <footer className="mt-12 text-center text-xs text-[#aaaaaa]">
        <p>© 2026 YSM 模型站. 基于 Mojang OreUI 设计语言.</p>
      </footer>
    </div>
  );
};

export default Login;
