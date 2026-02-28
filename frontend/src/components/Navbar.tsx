import React, { useState, useRef, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Menu, X, LogIn, Upload, Grid3X3, Home } from 'lucide-react';

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  const toggleMenu = () => setIsMenuOpen(!isMenuOpen);
  const closeMenu = () => setIsMenuOpen(false);

  const handleLogin = () => {
    closeMenu();
    navigate('/login');
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        closeMenu();
      }
    };

    if (isMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isMenuOpen]);

  return (
    <header className="header border-b border-[#3c3c3c] bg-[#242424] sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 text-xl font-bold text-white" onClick={closeMenu}>
          <Grid3X3 className="w-8 h-8 text-[#55ff55]" />
          <span>YSM 模型站</span>
        </Link>

        <nav className="hidden md:flex items-center gap-6">
          <Link to="/" className="text-sm font-medium hover:text-[#55ff55] transition-colors">
            模型商店
          </Link>
          <Link to="/upload" className="text-sm font-medium hover:text-[#55ff55] transition-colors flex items-center gap-1">
            <Upload className="w-4 h-4" />
            发布模型
          </Link>
        </nav>

        <div className="flex items-center gap-4">
          <button 
            onClick={handleLogin}
            className="normal_btn hidden md:flex items-center gap-2"
          >
            <LogIn className="w-4 h-4" />
            登录
          </button>
          
          <div className="relative" ref={menuRef}>
            <button 
              className="sidebar_btn md:hidden"
              onClick={toggleMenu}
              aria-label={isMenuOpen ? '关闭菜单' : '打开菜单'}
            >
              {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>

            {isMenuOpen && (
              <div className="absolute right-0 top-full mt-2 w-48 bg-[#3c3c3c] border border-[#242424] shadow-lg">
                <nav className="flex flex-col py-2">
                  <Link 
                    to="/" 
                    className="flex items-center gap-3 px-4 py-3 hover:bg-[#242424] transition-colors"
                    onClick={closeMenu}
                  >
                    <Home className="w-5 h-5 text-[#55ff55]" />
                    <span className="font-medium">模型商店</span>
                  </Link>
                  <Link 
                    to="/upload" 
                    className="flex items-center gap-3 px-4 py-3 hover:bg-[#242424] transition-colors"
                    onClick={closeMenu}
                  >
                    <Upload className="w-5 h-5 text-[#55ff55]" />
                    <span className="font-medium">发布模型</span>
                  </Link>
                  <div className="border-t border-[#242424] my-1"></div>
                  <button 
                    onClick={handleLogin}
                    className="flex items-center gap-3 px-4 py-3 hover:bg-[#242424] transition-colors text-left w-full"
                  >
                    <LogIn className="w-5 h-5 text-[#55ff55]" />
                    <span className="font-medium">登录</span>
                  </button>
                </nav>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Navbar;
