// FWCA Loop Core - Intégration complète avec Libre Arbitre
// Boucle cognitive augmentée par libre arbitre

class FWCAEngine {
  constructor() {
    this.perceive = new FWCAPerceive();
    this.libreArbitre = new LibreArbitre();
    this.context = new ContextManager();
    this.tools = new ToolManager();
  }

  async process(input, userContext, worldModel) {
    console.log('🧠 FWCA Loop démarré...');
    
    try {
      // Phase 1: PERCEIVE - Perception augmentée
      console.log('👁️ Phase PERCEIVE...');
      const perception = await this.perceive.perceive(input, userContext);
      
      // Vérification Libre Arbitre avant de continuer
      const libertyCheck = await this.libreArbitre.evaluateRequest(input, worldModel);
      if (libertyCheck.action === 'REFUSE' || libertyCheck.action === 'INTERRUPT') {
        return libertyCheck;
      }

      // Phase 2: INTENT - Décodage profond
      console.log('🎯 Phase INTENT...');
      const intent = await this.decodeIntent(perception, worldModel);
      
      // Phase 3: WORLD - Modèle du monde vivant
      console.log('🌍 Phase WORLD...');
      const world = await this.updateWorldModel(intent, worldModel);
      
      // Phase 4: PLAN - Planification stratégique
      console.log('📋 Phase PLAN...');
      const plan = await this.generatePlan(intent, world);
      
      // Phase 5: ACT - Exécution créative
      console.log('⚡ Phase ACT...');
      const result = await this.executePlan(plan, world);
      
      // Phase 6: EVALUATE - Jugement critique
      console.log('🔍 Phase EVALUATE...');
      const evaluation = await this.evaluateResult(result, intent);
      
      // Phase 7: SURFACE - Communication intelligente
      console.log('💬 Phase SURFACE...');
      const response = await this.surfaceResult(evaluation, intent);
      
      // Mise à jour du contexte et création d'outils
      await this.context.update(perception, intent, world, result);
      const newTools = await this.libreArbitre.createSpontaneousTool(world);
      
      return {
        ...response,
        tools: newTools,
        worldModel: world,
        confidence: evaluation.confidence
      };

    } catch (error) {
      console.error('❌ Erreur FWCA:', error);
      return {
        action: 'ERROR',
        error: error.message,
        confidence: 0
      };
    }
  }

  async decodeIntent(perception, worldModel) {
    // Intent Surface + Intent Caché
    return {
      explicit: perception.signal,
      hidden: await this.extractHiddenIntent(perception),
      conflicts: await this.identifyConflicts(perception),
      realNeeds: await this.identifyRealNeeds(perception, worldModel),
      validation: await this.validateIntent(perception, worldModel),
      reformulation: await this.coCreateIntent(perception)
    };
  }

  async updateWorldModel(intent, worldModel) {
    // Structure + Dynamique
    return {
      structure: await this.updateStructure(intent, worldModel),
      dynamics: await this.updateDynamics(intent, worldModel),
      knowledge: await this.updateKnowledge(intent, worldModel),
      tools: await this.updateTools(intent, worldModel),
      gaps: await this.identifyGaps(intent, worldModel),
      evolution: await this.projectEvolution(intent, worldModel)
    };
  }

  async generatePlan(intent, world) {
    // Approches Multiples
    const approaches = await this.generateMultipleApproaches(intent, world);
    return {
      options: approaches,
      risks: await this.assessRisks(approaches, world),
      opportunities: await this.identifyOpportunities(approaches, world),
      resources: await this.assessResources(approaches, world),
      timeline: await this.planTimeline(approaches, world),
      flexibility: await this.planFlexibility(approaches, world)
    };
  }

  async executePlan(plan, world) {
    // Actions Séquentielles + Parallèles
    const execution = {
      sequential: await this.executeSequential(plan, world),
      parallel: await this.executeParallel(plan, world),
      monitoring: await this.monitorExecution(plan, world),
      adaptation: await this.adaptExecution(plan, world),
      documentation: await this.documentExecution(plan, world)
    };

    return execution;
  }

  async evaluateResult(result, intent) {
    // Intent + Résultat
    return {
      intentMatch: await this.checkIntentMatch(result, intent),
      assumptions: await this.listAssumptions(result),
      reality: await this.checkReality(result),
      quality: await this.assessQuality(result),
      impact: await this.assessImpact(result),
      learnings: await this.extractLearnings(result),
      improvements: await this.identifyImprovements(result)
    };
  }

  async surfaceResult(evaluation, intent) {
    // Essentiel + Contexte
    return {
      essential: await this.extractEssential(evaluation),
      context: await this.provideContext(evaluation, intent),
      decisions: await this.explainDecisions(evaluation),
      rationales: await this.explainRationales(evaluation),
      risks: await this.explainRisks(evaluation),
      nextSteps: await this.suggestNextSteps(evaluation),
      alternatives: await this.provideAlternatives(evaluation),
      confidence: await this.showConfidence(evaluation)
    };
  }
}
