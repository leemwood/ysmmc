import React from 'react';
import { Link } from 'react-router-dom';
import { Home, Search, Grid3X3 } from 'lucide-react';

const NotFound: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#242424] text-white flex flex-col items-center justify-center p-4">
      <div className="text-center">
        <Grid3X3 className="w-24 h-24 text-[#55ff55] mx-auto mb-8 opacity-50" />
        <h1 className="text-6xl font-bold mb-4">404</h1>
        <h2 className="text-2xl font-bold mb-4">页面未找到</h2>
        <p className="text-[#aaaaaa] mb-8 max-w-md">
          您访问的页面不存在或已被移除。请检查网址是否正确，或返回首页继续浏览。
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link 
            to="/" 
            className="green_btn px-8 py-3 flex items-center justify-center gap-2 font-bold"
          >
            <Home className="w-5 h-5" />
            返回首页
          </Link>
          <Link 
            to="/" 
            className="normal_btn px-8 py-3 flex items-center justify-center gap-2 font-bold"
          >
            <Search className="w-5 h-5" />
            浏览模型
          </Link>
        </div>
      </div>
    </div>
  );
};

export default NotFound;
