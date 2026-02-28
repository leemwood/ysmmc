import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { UserPlus, Mail, Lock, User, AlertCircle, Loader2, Grid3X3 } from 'lucide-react';
import api from '../api/client';
import type { ApiResponse } from '../types';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (password !== confirmPassword) {
      setError('两次输入的密码不一致');
      return;
    }

    setLoading(true);

    try {
      const response = await api.post<ApiResponse<any>>('/auth/register', {
        username,
        email,
        password,
      });

      if (response.data.success) {
        navigate('/login');
      } else {
        setError(response.data.message || '注册失败');
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
            <UserPlus className="w-8 h-8 text-[#55ff55]" />
          </div>
          <h2 className="text-2xl font-bold tracking-tight">加入我们</h2>
          <p className="text-sm text-[#aaaaaa] mt-2">创建一个账户以开始分享您的模型</p>
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
              <User className="w-3.5 h-3.5" />
              用户名
            </label>
            <input
              type="text"
              required
              placeholder="您的称呼"
              className="input w-full bg-[#242424] border-[#3c3c3c] h-12 px-4 focus:border-[#55ff55] transition-all"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>

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
            <label className="text-xs font-bold uppercase tracking-wider text-[#aaaaaa] flex items-center gap-2">
              <Lock className="w-3.5 h-3.5" />
              密码
            </label>
            <input
              type="password"
              required
              placeholder="••••••••"
              className="input w-full bg-[#242424] border-[#3c3c3c] h-12 px-4 focus:border-[#55ff55] transition-all"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label className="text-xs font-bold uppercase tracking-wider text-[#aaaaaa] flex items-center gap-2">
              <Lock className="w-3.5 h-3.5" />
              确认密码
            </label>
            <input
              type="password"
              required
              placeholder="••••••••"
              className="input w-full bg-[#242424] border-[#3c3c3c] h-12 px-4 focus:border-[#55ff55] transition-all"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
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
                注册中...
              </>
            ) : (
              '创建账户'
            )}
          </button>
        </form>

        <div className="mt-8 pt-6 border-t border-[#242424] text-center">
          <p className="text-sm text-[#aaaaaa]">
            已有账号?{' '}
            <Link to="/login" className="text-[#55ff55] font-bold hover:underline">
              立即登录
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

export default Register;
