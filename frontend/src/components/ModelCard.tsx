import React from 'react';
import type { Model } from '../types';
import { Download, Heart, Eye } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface ModelCardProps {
  model: Model;
}

const ModelCard: React.FC<ModelCardProps> = ({ model }) => {
  const navigate = useNavigate();

  return (
    <div 
      className="card hover:border-[#55ff55] transition-all cursor-pointer group bg-[#3c3c3c] border border-[#242424] overflow-hidden"
      onClick={() => navigate(`/model/${model.id}`)}
    >
      <div className="relative aspect-video overflow-hidden bg-[#242424]">
        <img 
          src={model.image_url || '/placeholder-model.png'} 
          alt={model.title}
          className="w-full h-full object-cover transition-transform group-hover:scale-105"
        />
        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-4">
          <div className="flex flex-col items-center text-white">
            <Eye className="w-6 h-6 mb-1" />
            <span className="text-xs font-medium">查看详情</span>
          </div>
        </div>
      </div>
      
      <div className="p-4 flex flex-col gap-3">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-bold text-lg text-white line-clamp-1 leading-tight">{model.title}</h3>
          <span className="badge text-[10px] px-2 py-0.5 bg-[#55ff55] text-black uppercase font-bold tracking-wider">
            {model.tags[0] || '模型'}
          </span>
        </div>
        
        <p className="text-xs text-[#aaaaaa] line-clamp-2 h-8 leading-relaxed">
          {model.description || '该模型暂无描述信息...'}
        </p>

        <div className="flex items-center justify-between pt-3 border-t border-[#242424]">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1.5 text-xs text-[#aaaaaa]">
              <Download className="w-3.5 h-3.5" />
              <span>{model.downloads}</span>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-[#aaaaaa]">
              <Heart className="w-3.5 h-3.5" />
              <span>0</span>
            </div>
          </div>
          <button 
            className="normal_btn text-[10px] py-1 px-3 min-h-0 h-7 flex items-center gap-1"
            onClick={(e) => {
              e.stopPropagation();
              // TODO: Handle download
            }}
          >
            下载
          </button>
        </div>
      </div>
    </div>
  );
};

export default ModelCard;
