// Gestion de Contexte Avancée - Multi-Niveaux avec Conscience
// Prometheus sait CE qu'il sait et CE qu'il ignore

class ContextManager {
  constructor() {
    this.levels = {
      session: [],      // Contexte immédiat de la conversation
      world: {},       // Contexte profond du monde actuel  
      universal: {},   // Contexte à travers tous les mondes
      temporal: []     // Évolution temporelle des contextes
    };
    this.compression = new ContextCompression();
    this.retrieval = new ContextRetrieval();
  }

  async update(perception, intent, world, result) {
    // Cross-Context Sync
    await this.updateSession(perception, intent, result);
    await this.updateWorld(world);
    await this.updateUniversal(world);
    await this.updateTemporal(perception, intent, world, result);
    
    // Compression automatique si nécessaire
    if (await this.needsCompression()) {
      await this.compression.compress(this);
    }
  }

  async updateSession(perception, intent, result) {
    // Session Context - derniers 10 échanges
    this.levels.session.push({
      timestamp: Date.now(),
      perception: perception,
      intent: intent,
      result: result
    });

    // Garder seulement les 10 derniers
    if (this.levels.session.length > 10) {
      this.levels.session.shift();
    }
  }

  async updateWorld(world) {
    // World Context - état actuel du monde
    this.levels.world = {
      model: world.structure,
      tools: world.tools,
      knowledge: world.knowledge,
      gaps: world.gaps,
      evolution: world.evolution
    };
  }

  async updateUniversal(world) {
    // Universal Context - patterns cross-mondes
    if (!this.levels.worlds) {
      this.levels.worlds = [];
    }
    
    this.levels.worlds.push({
      timestamp: Date.now(),
      world: world
    });

    // Extraire patterns universels
    this.levels.universal = await this.extractUniversalPatterns(this.levels.worlds);
  }

  async updateTemporal(perception, intent, world, result) {
    // Temporal Context - évolution temporelle
    this.levels.temporal.push({
      timestamp: Date.now(),
      event: 'interaction',
      perception: perception,
      intent: intent,
      world: world,
      result: result
    });

    // Analyser les tendances temporelles
    this.levels.temporal.trends = await this.analyzeTemporalTrends(this.levels.temporal);
  }

  async search(query, maxResults = 5) {
    // Semantic Search avec relevance scoring
    const results = await this.retrieval.semanticSearch(query, this, maxResults);
    return results;
  }

  async getContextForPrompt(currentIntent) {
    // Context Injection - injecter contexte pertinent
    const relevantContext = {
      recent: await this.getRecentContext(3),
      world: await this.getWorldContext(currentIntent),
      patterns: await this.getPatternContext(currentIntent),
      temporal: await this.getTemporalContext(currentIntent)
    };

    return await this.formatContextForPrompt(relevantContext);
  }

  async detectConflicts() {
    // Conflict Resolution - gérer contradictions contextuelles
    const conflicts = [];
    
    // Vérifier cohérence entre niveaux
    const sessionWorld = await this.checkSessionWorldConsistency();
    if (sessionWorld.conflicts.length > 0) {
      conflicts.push(...sessionWorld.conflicts);
    }

    const worldUniversal = await this.checkWorldUniversalConsistency();
    if (worldUniversal.conflicts.length > 0) {
      conflicts.push(...worldUniversal.conflicts);
    }

    return conflicts;
  }

  async getMetaContext() {
    // Méta-contexte - conscience de ce qu'on sait/ignore
    return {
      known: await this.summarizeKnown(),
      unknown: await this.identifyUnknown(),
      confidence: await this.calculateContextConfidence(),
      gaps: await this.identifyKnowledgeGaps(),
      evolution: await this.getContextEvolution()
    };
  }
}

// Context Compression Algorithm
class ContextCompression {
  async compress(contextManager) {
    // Information Hierarchy
    const hierarchy = await this.identifyInformationHierarchy(contextManager);
    
    // Semantic Clustering
    const clusters = await this.clusterSemanticInformation(contextManager);
    
    // Temporal Decay
    const decayed = await this.applyTemporalDecay(contextManager);
    
    // Pattern Extraction
    const patterns = await this.extractPatterns(contextManager);
    
    // Size Optimization - maintenir < 4K tokens
    const compressed = await this.optimizeSize(hierarchy, clusters, decayed, patterns);
    
    return compressed;
  }

  async identifyInformationHierarchy(contextManager) {
    // Identifier informations critiques vs secondaires
    return {
      critical: await this.extractCritical(contextManager),
      important: await this.extractImportant(contextManager),
      secondary: await this.extractSecondary(contextManager)
    };
  }

  async clusterSemanticInformation(contextManager) {
    // Grouper informations par sémantique
    const allInfo = await this.getAllInformation(contextManager);
    const clusters = await this.semanticCluster(allInfo);
    
    return clusters;
  }

  async applyTemporalDecay(contextManager) {
    // Réduire importance d'anciennes informations
    const now = Date.now();
    const decayed = {};
    
    for (const [level, data] of Object.entries(contextManager.levels)) {
      decayed[level] = await this.applyDecayToLevel(data, now);
    }
    
    return decayed;
  }

  async extractPatterns(contextManager) {
    // Extraire patterns récurrents
    return {
      temporal: await this.extractTemporalPatterns(contextManager),
      semantic: await this.extractSemanticPatterns(contextManager),
      structural: await this.extractStructuralPatterns(contextManager)
    };
  }
}

// Context Retrieval System
class ContextRetrieval {
  async semanticSearch(query, contextManager, maxResults) {
    // Rechercher par sens pas par mots-clés
    const queryEmbedding = await this.generateEmbedding(query);
    const candidates = await this.findCandidates(queryEmbedding, contextManager);
    const scored = await this.scoreRelevance(candidates, queryEmbedding);
    
    return scored.slice(0, maxResults);
  }

  async scoreRelevance(candidates, queryEmbedding) {
    // Noter pertinence des résultats
    const scored = [];
    
    for (const candidate of candidates) {
      const candidateEmbedding = await this.generateEmbedding(candidate);
      const similarity = await this.calculateSimilarity(queryEmbedding, candidateEmbedding);
      const relevance = await this.calculateRelevance(candidate, queryEmbedding);
      
      scored.push({
        ...candidate,
        similarity,
        relevance,
        score: similarity * relevance
      });
    }
    
    return scored.sort((a, b) => b.score - a.score);
  }

  async dynamicUpdate(contextManager, newInfo) {
    // Dynamic Updating - mettre à jour contexte en temps réel
    await this.updateIndexes(contextManager, newInfo);
    await this.refreshEmbeddings(contextManager, newInfo);
    await this.recalculateRelevance(contextManager, newInfo);
  }
}
