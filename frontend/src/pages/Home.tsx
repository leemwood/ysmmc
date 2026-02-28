import React, { useEffect, useState } from 'react';
import Navbar from '../components/Navbar';
import ModelCard from '../components/ModelCard';
import { modelApi, announcementApi } from '../api';
import type { Model, Announcement } from '../types';
import { Search, Megaphone, Loader2, Sparkles, AlertCircle, Grid3X3 } from 'lucide-react';

const Home: React.FC = () => {
  const [models, setModels] = useState<Model[]>([]);
  const [announcements, setAnnouncements] = useState<Announcement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [modelsRes, announcementsRes] = await Promise.all([
          modelApi.list(1, 12, searchQuery),
          announcementApi.list(),
        ]);
        setModels(modelsRes.data.data.items);
        setAnnouncements(announcementsRes.data.data);
      } catch (err) {
        setError('无法连接到 API 服务器，请稍后再试。');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [searchQuery]);

  return (
    <div className="min-h-screen bg-[#242424] text-white flex flex-col">
      <Navbar />

      <main className="flex-1 container mx-auto px-4 py-8">
        {/* Banner Section */}
        <section className="mb-12 relative overflow-hidden bg-gradient-to-r from-[#1a1a1a] to-[#333333] border border-[#3c3c3c] rounded-lg p-10">
          <div className="relative z-10 max-w-2xl">
            <h1 className="text-4xl font-extrabold mb-4 tracking-tight leading-tight flex items-center gap-3">
              <Sparkles className="w-8 h-8 text-[#55ff55]" />
              探索无限可能的 YSM 模型
            </h1>
            <p className="text-lg text-[#aaaaaa] mb-8 leading-relaxed">
              这里是 YSM 模型站，为您提供最全面、最优质的模型资源，助力您的创作之旅。
            </p>
            <div className="flex items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[#aaaaaa]" />
                <input 
                  type="text" 
                  placeholder="搜索您感兴趣的模型..."
                  className="input w-full pl-12 pr-4 h-12 bg-[#242424] border-[#3c3c3c] focus:border-[#55ff55] transition-all"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>
              <button className="green_btn h-12 px-8 flex items-center gap-2 font-bold uppercase tracking-wider">
                搜索
              </button>
            </div>
          </div>
          <div className="absolute right-0 top-0 w-1/3 h-full opacity-10 flex items-center justify-center pointer-events-none">
            <Grid3X3 className="w-64 h-64" />
          </div>
        </section>

        {/* Announcements Section */}
        {announcements.length > 0 && (
          <section className="mb-12">
            <div className="flex items-center gap-2 mb-6 border-l-4 border-[#55ff55] pl-4">
              <Megaphone className="w-6 h-6 text-[#55ff55]" />
              <h2 className="text-2xl font-bold tracking-tight">最新公告</h2>
            </div>
            <div className="grid gap-4">
              {announcements.map((ann) => (
                <div key={ann.id} className="card bg-[#3c3c3c] border border-[#242424] p-5 flex items-start gap-4 hover:border-[#55ff55] transition-all group cursor-pointer">
                  <div className="w-10 h-10 rounded-full bg-[#242424] flex items-center justify-center shrink-0 border border-[#3c3c3c] group-hover:bg-[#55ff55] group-hover:text-black transition-all">
                    <Megaphone className="w-5 h-5" />
                  </div>
                  <div>
                    <h3 className="font-bold text-lg mb-1">{ann.title}</h3>
                    <p className="text-sm text-[#aaaaaa] leading-relaxed">{ann.content}</p>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}

        {/* Models Grid */}
        <section>
          <div className="flex items-center justify-between mb-8 border-b border-[#3c3c3c] pb-4">
            <div className="flex items-center gap-2">
              <Grid3X3 className="w-6 h-6 text-[#55ff55]" />
              <h2 className="text-2xl font-bold tracking-tight">模型商店</h2>
            </div>
            <div className="flex items-center gap-4 text-xs font-medium text-[#aaaaaa]">
              <span>共 {models.length} 个模型</span>
            </div>
          </div>

          {loading ? (
            <div className="flex flex-col items-center justify-center py-32 gap-4">
              <Loader2 className="w-12 h-12 text-[#55ff55] animate-spin" />
              <p className="text-sm font-medium animate-pulse">正在加载模型商店...</p>
            </div>
          ) : error ? (
            <div className="flex flex-col items-center justify-center py-32 gap-4 border-2 border-dashed border-red-500/20 rounded-lg bg-red-500/5">
              <AlertCircle className="w-12 h-12 text-red-500" />
              <div className="text-center">
                <h3 className="text-xl font-bold mb-2">出错了</h3>
                <p className="text-[#aaaaaa] mb-6">{error}</p>
                <button 
                  onClick={() => window.location.reload()}
                  className="red_btn px-8 py-2 font-bold"
                >
                  重新加载
                </button>
              </div>
            </div>
          ) : models.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-32 gap-4 border-2 border-dashed border-[#3c3c3c] rounded-lg">
              <Search className="w-12 h-12 text-[#aaaaaa]" />
              <p className="text-[#aaaaaa] font-medium text-lg">未找到相关模型...</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
              {models.map((model) => (
                <ModelCard key={model.id} model={model} />
              ))}
            </div>
          )}
        </section>
      </main>

      <footer className="mt-auto border-t border-[#3c3c3c] py-10 bg-[#1a1a1a]">
        <div className="container mx-auto px-4 text-center">
          <p className="text-sm text-[#aaaaaa]">
            © 2026 YSM 模型站. 基于 Mojang OreUI 设计语言构建.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Home;
